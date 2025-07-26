package media

import (
	"bytes"
	"fmt"
	"image"

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
	var buffer bytes.Buffer

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

	if color {
		// ANSI 컬러 모드
		for y := 0; y < newHeight; y++ {
			for x := 0; x < width; x++ {
				vec := resized.GetVecbAt(y, x) // BGR 형식
				b, g, r := vec[0], vec[1], vec[2]

				// 문자 선택 (컬러 강조용 단순화)
				var ch rune
				if r > g && r > b {
					ch = '@' // 빨강 계열
				} else if g > r && g > b {
					ch = '&' // 초록 계열
				} else if b > r && b > g {
					ch = '#' // 파랑 계열
				} else {
					ch = '.' // 기타
				}

				// ANSI 24bit 컬러 적용
				buffer.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c", r, g, b, ch))
			}
			buffer.WriteString("\x1b[0m\n") // 한 줄 끝나면 색상 리셋 + 줄바꿈
		}
	} else {
		// 흑백 ASCII 모드
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)
		for y := 0; y < newHeight; y++ {
			for x := 0; x < width; x++ {
				pixel := gray.GetUCharAt(y, x)
				idx := int(float64(pixel) / 255.0 * float64(len(c.charset)))
				if idx >= len(c.charset) {
					idx = len(c.charset) - 1
				}
				buffer.WriteRune(c.charset[idx])
			}
			buffer.WriteRune('\n')
		}
	}

	return buffer.String(), nil
}
