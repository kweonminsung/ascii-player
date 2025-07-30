package media

import (
	"fmt"
	"image"

	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)

type AnsiConverter struct{}

// NewAnsiConverter는 AnsiConverter 인스턴스를 생성합니다.
func NewAnsiConverter() *AnsiConverter {
	return &AnsiConverter{}
}

// ConvertToScreen은 gocv.Mat을 ANSI로 변환하여 tcell.Screen에 직접 렌더링합니다.
func (c *AnsiConverter) ConvertToScreen(img gocv.Mat, width, height int, color bool, screen tcell.Screen) error {
	originalWidth := float64(img.Cols())
	originalHeight := float64(img.Rows())
	if originalWidth == 0 || originalHeight == 0 {
		return fmt.Errorf("invalid image dimensions: %dx%d", int(originalWidth), int(originalHeight))
	}

	newWidth, newHeight := width, height

	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Pt(newWidth, newHeight), 0, 0, gocv.InterpolationLinear)

	if !color {
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)
		resized = gray
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			style := tcell.StyleDefault
			if color {
				vec := resized.GetVecbAt(y, x)
				b, g, r := vec[0], vec[1], vec[2]
				style = style.Foreground(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			} else {
				pixel := resized.GetUCharAt(y, x)
				style = style.Foreground(tcell.NewRGBColor(int32(pixel), int32(pixel), int32(pixel)))
			}
			screen.SetContent(x, y, '█', nil, style)
		}
	}

	return nil
}
