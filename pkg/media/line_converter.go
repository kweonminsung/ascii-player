package media

import (
	"bytes"
	"image"
	"math"

	"github.com/kweonminsung/console-cinema/pkg/types"
	"gocv.io/x/gocv"
)

// LineCharset은 라인 렌더링에 사용될 문자 집합입니다.
const LineCharset = `|/—\`

// LineConverter는 이미지의 경계선을 감지하여 라인 문자로 변환하는 기능을 제공합니다.
type LineConverter struct {
	charset           []rune
	gradientThreshold float64
}

// NewLineConverter는 새로운 LineConverter 인스턴스를 생성합니다.
func NewLineConverter(charset string, threshold float64) *LineConverter {
	if charset == "" {
		charset = LineCharset
	}
	if threshold <= 0 {
		threshold = 30.0 // 기본 임계값
	}
	return &LineConverter{
		charset:           []rune(charset),
		gradientThreshold: threshold,
	}
}

// Convert는 gocv.Mat 이미지를 경계선 기반의 ASCII 문자로 변환합니다.
func (c *LineConverter) Convert(img gocv.Mat, width, height int) (string, error) {
	var buffer bytes.Buffer

	// 1. 이미지 비율에 맞게 높이 재계산
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

	// 2. 이미지 크기 조절 및 흑백 변환
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: width, Y: newHeight}, 0, 0, gocv.InterpolationLinear)

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

	// 3. Sobel 필터를 사용하여 그래디언트 계산
	gradX := gocv.NewMat()
	gradY := gocv.NewMat()
	defer gradX.Close()
	defer gradY.Close()

	gocv.Sobel(gray, &gradX, gocv.MatTypeCV16S, 1, 0, 3, 1, 0, gocv.BorderDefault)
	gocv.Sobel(gray, &gradY, gocv.MatTypeCV16S, 0, 1, 3, 1, 0, gocv.BorderDefault)

	// 4. 픽셀을 순회하며 각도에 따라 문자로 변환
	for y := 0; y < newHeight; y++ {
		for x := 0; x < width; x++ {
			dx := float64(gradX.GetShortAt(y, x))
			dy := float64(gradY.GetShortAt(y, x))

			// 그래디언트 크기 계산
			magnitude := math.Sqrt(dx*dx + dy*dy)

			if magnitude < c.gradientThreshold {
				buffer.WriteRune(' ') // 임계값보다 작으면 공백 처리
				continue
			}

			// 각도 계산 (라디안)
			angle := math.Atan2(dy, dx)
			// 각도를 0-180도로 변환
			angleDegrees := angle * (180.0 / math.Pi)
			if angleDegrees < 0 {
				angleDegrees += 180
			}

			// 각도에 따라 문자 선택
			var char rune
			switch {
			case (angleDegrees >= 0 && angleDegrees < 22.5) || (angleDegrees >= 157.5 && angleDegrees <= 180):
				char = '—' // 수평선
			case angleDegrees >= 22.5 && angleDegrees < 67.5:
				char = '/' // 대각선
			case angleDegrees >= 67.5 && angleDegrees < 112.5:
				char = '|' // 수직선
			case angleDegrees >= 112.5 && angleDegrees < 157.5:
				char = '\\' // 역대각선
			default:
				char = ' '
			}
			buffer.WriteRune(char)
		}
		buffer.WriteRune('\n')
	}

	return buffer.String(), nil
}
