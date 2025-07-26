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

type AnsiConverter struct{}

// NewAnsiConverter는 AnsiConverter 인스턴스를 생성합니다.
func NewAnsiConverter() *AnsiConverter {
	return &AnsiConverter{}
}

// Render는 gocv.Mat을 ANSI 컬러 문자열(픽셀당 컬러 문자)로 변환하여 반환합니다.
// color가 true일 경우 ANSI 색상 출력, false일 경우 회색조
func (c *AnsiConverter) Convert(img gocv.Mat, width, height int, color bool) (string, error) {
	originalWidth := float64(img.Cols())
	originalHeight := float64(img.Rows())
	if originalWidth == 0 {
		return "", fmt.Errorf("invalid image width: 0")
	}
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
	gocv.Resize(img, &resized, image.Pt(width, newHeight), 0, 0, gocv.InterpolationLinear)

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
						vec := resized.GetVecbAt(y, x) // BGR
						b, g, r := vec[0], vec[1], vec[2]
						ch := '█'
						line.WriteString(fmt.Sprintf("[#%02x%02x%02x]%c", r, g, b, ch))
					}
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
						val := gray.GetUCharAt(y, x)
						ch := '█'
						line.WriteString(fmt.Sprintf("[#%02x%02x%02x]%c", val, val, val, ch))
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
		buffer.WriteString("\n")
	}

	return buffer.String(), nil
}
