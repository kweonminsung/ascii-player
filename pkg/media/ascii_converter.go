package media

import (
	"bytes"
	"fmt"
	"image"
	"runtime"
	"sync"

	"github.com/kweonminsung/ascii-player/pkg/types"
	"gocv.io/x/gocv"
)

// DefaultCharset은 밝기 순서에 따라 정렬된 기본 ASCII 문자 집합입니다.
const DefaultCharset = " .:-=+*#%@"

// AsciiConverter는 이미지를 ASCII 아트로 변환하는 기능을 제공합니다.
type AsciiConverter struct {
	charset []rune
}

// NewAsciiConverter는 새로운 AsciiConverter 인스턴스를 생성합니다.
func NewAsciiConverter() *AsciiConverter {
	return &AsciiConverter{
		charset: []rune(DefaultCharset),
	}
}

// Convert는 gocv.Mat 이미지를 지정된 너비와 높이의 ASCII 문자열로 변환합니다.
// color가 true이면 ANSI 컬러 코드 기반 ASCII 출력, false이면 흑백 ASCII 출력.
func (c *AsciiConverter) Convert(img gocv.Mat, width, height int, color bool) (string, error) {
	originalWidth := float64(img.Cols())
	originalHeight := float64(img.Rows())
	aspectRatio := originalHeight / originalWidth
	newHeight := int(float64(width) * aspectRatio * types.YScaleFactor)
	if newHeight <= 0 {
		newHeight = 1
	}
	if height > 0 && newHeight > height {
		newHeight = height
	}

	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: width, Y: newHeight}, 0, 0, gocv.InterpolationLinear)

	lines := make([]string, newHeight)
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	if color {
		rowJobs := make(chan int, newHeight)
		for y := 0; y < newHeight; y++ {
			rowJobs <- y
		}
		close(rowJobs)

		wg.Add(newHeight)
		for i := 0; i < numWorkers; i++ {
			go func() {
				for y := range rowJobs {
					var line bytes.Buffer
					for x := 0; x < width; x++ {
						vec := resized.GetVecbAt(y, x) // BGR 형식
						b, g, r := vec[0], vec[1], vec[2]
						var ch rune
						if r > g && r > b {
							ch = '@'
						} else if g > r && g > b {
							ch = '&'
						} else if b > r && b > g {
							ch = '#'
						} else {
							ch = '.'
						}
						line.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c", r, g, b, ch))
					}
					line.WriteString("\x1b[0m")
					lines[y] = line.String()
					wg.Done()
				}
			}()
		}
		wg.Wait()
	} else {
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

		rowJobs := make(chan int, newHeight)
		for y := 0; y < newHeight; y++ {
			rowJobs <- y
		}
		close(rowJobs)

		wg.Add(newHeight)
		for i := 0; i < numWorkers; i++ {
			go func() {
				for y := range rowJobs {
					var line bytes.Buffer
					for x := 0; x < width; x++ {
						pixel := gray.GetUCharAt(y, x)
						idx := int(float64(pixel) / 255.0 * float64(len(c.charset)))
						if idx >= len(c.charset) {
							idx = len(c.charset) - 1
						}
						line.WriteRune(c.charset[idx])
					}
					lines[y] = line.String()
					wg.Done()
				}
			}()
		}
		wg.Wait()
	}

	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line)
		buffer.WriteRune('\n')
	}

	return buffer.String(), nil
}
