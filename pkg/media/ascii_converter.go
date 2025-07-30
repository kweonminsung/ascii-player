package media

import (
	"fmt"
	"image"

	"github.com/gdamore/tcell/v2"
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

// ConvertToScreen은 gocv.Mat 이미지를 ASCII로 변환하여 tcell.Screen에 직접 렌더링합니다.
func (c *AsciiConverter) ConvertToScreen(img gocv.Mat, width, height int, color bool, screen tcell.Screen) error {
	originalWidth := float64(img.Cols())
	originalHeight := float64(img.Rows())
	if originalWidth == 0 || originalHeight == 0 {
		return fmt.Errorf("invalid image dimensions: %dx%d", int(originalWidth), int(originalHeight))
	}

	newWidth, newHeight := width, height

	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationLinear)

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			pixel := gray.GetUCharAt(y, x)
			idx := int(float64(pixel) / 255.0 * float64(len(c.charset)))
			if idx >= len(c.charset) {
				idx = len(c.charset) - 1
			}
			ch := c.charset[idx]

			style := tcell.StyleDefault
			if color {
				vec := resized.GetVecbAt(y, x)
				b, g, r := vec[0], vec[1], vec[2]
				style = style.Foreground(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			}
			screen.SetContent(x, y, ch, nil, style)
		}
	}

	return nil
}
