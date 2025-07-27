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

	var buffer = make([][]byte, newHeight)
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	rowJobs := make(chan int, newHeight)
	for y := 0; y < newHeight; y++ {
		rowJobs <- y
	}
	close(rowJobs)

	if color {
		data := resized.ToBytes()
		lineSize := width * 3

		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func() {
				defer wg.Done()
				for y := range rowJobs {
					offset := y * lineSize
					var line []byte
					prevColor := ""

					for x := 0; x < width; x++ {
						b := data[offset+3*x]
						g := data[offset+3*x+1]
						r := data[offset+3*x+2]

						colorTag := fmt.Sprintf("%02x%02x%02x", r, g, b)
						if colorTag != prevColor {
							line = append(line, '[')
							line = append(line, '#')
							line = append(line, colorTag...)
							line = append(line, ']')
							prevColor = colorTag
						}

						line = append(line, []byte("█")...)
					}
					buffer[y] = append(line, '\n')
				}
			}()
		}
	} else {
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)
		data := gray.ToBytes()

		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func() {
				defer wg.Done()
				for y := range rowJobs {
					offset := y * width
					var line []byte
					prevGray := -1

					for x := 0; x < width; x++ {
						val := int(data[offset+x])

						if val != prevGray {
							line = append(line, '[')
							line = append(line, '#')
							hex := fmt.Sprintf("%02x%02x%02x", val, val, val)
							line = append(line, hex...)
							line = append(line, ']')
							prevGray = val
						}

						line = append(line, []byte("█")...)
					}
					buffer[y] = append(line, '\n')
				}
			}()
		}
	}

	wg.Wait()

	var final bytes.Buffer
	for _, line := range buffer {
		final.Write(line)
	}
	return final.String(), nil
}
