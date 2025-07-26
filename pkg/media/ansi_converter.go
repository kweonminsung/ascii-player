package media

import (
	"bytes"
	"fmt"
	"image"

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
	var buffer bytes.Buffer

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

	if color {
		// 컬러 ANSI 모드
		for y := 0; y < newHeight; y++ {
			for x := 0; x < width; x++ {
				vec := resized.GetVecbAt(y, x) // BGR
				b, g, r := vec[0], vec[1], vec[2]

				// 문자 선택
				ch := '█'

				// ANSI 24bit 전경색 설정 + 문자 출력
				buffer.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c", r, g, b, ch))
			}
			buffer.WriteString("\x1b[0m\n")
		}
	} else {
		// 흑백 ANSI 모드 (회색조)
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

		for y := 0; y < newHeight; y++ {
			for x := 0; x < width; x++ {
				val := gray.GetUCharAt(y, x)
				ch := '█' // 밝기 기반 문자는 개선 가능
				// foreground 색상: 회색
				buffer.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c", val, val, val, ch))
			}
			buffer.WriteString("\x1b[0m\n")
		}
	}

	return buffer.String(), nil
}
