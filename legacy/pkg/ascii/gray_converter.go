package ascii

import (
	"bytes"
	"image"

	"gocv.io/x/gocv"
)

// DefaultCharset은 ASCII 변환에 사용될 기본 문자 집합입니다.
// 밝기 순서에 따라 문자를 정렬합니다.
const DefaultCharset = " .:-=+*#%@"

// YScaleFactor는 터미널 문자의 세로/가로 비율을 보정하기 위한 값입니다.
// 일반적인 터미널 폰트는 높이가 너비의 약 2배이므로 0.5에 가까운 값을 사용합니다.
const YScaleFactor = 0.55

// GrayConverter는 이미지를 ASCII 아트로 변환하는 기능을 제공합니다.
type GrayConverter struct {
	charset []rune
}

// NewGrayConverter는 새로운 GrayConverter 인스턴스를 생성합니다.
// charset이 비어있으면 DefaultCharset을 사용합니다.
func NewGrayConverter(charset string) *GrayConverter {
	if charset == "" {
		charset = DefaultCharset
	}
	return &GrayConverter{
		charset: []rune(charset),
	}
}

// Convert는 gocv.Mat 이미지를 지정된 너비와 높이의 ASCII 문자열로 변환합니다.
// 터미널의 문자 비율에 맞게 이미지 높이를 자동으로 조절합니다.
func (c *GrayConverter) Convert(img gocv.Mat, width, height int) (string, error) {
	// 최종 ASCII 아트 출력을 위한 버퍼
	var buffer bytes.Buffer

	// 1. 이미지 비율에 맞게 높이 재계산
	originalWidth := float64(img.Cols())
	originalHeight := float64(img.Rows())
	aspectRatio := originalHeight / originalWidth
	newHeight := int(float64(width) * aspectRatio * YScaleFactor)
	if newHeight <= 0 {
		newHeight = 1 // 높이는 최소 1 이상이어야 합니다.
	}
	// 사용자가 지정한 height와 계산된 newHeight 중 더 작은 값을 사용 (선택적)
	if height > 0 && newHeight > height {
		newHeight = height
	}


	// 2. 이미지 크기 조절
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: width, Y: newHeight}, 0, 0, gocv.InterpolationLinear)

	// 3. 흑백으로 변환
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

	// 4. 픽셀을 순회하며 ASCII 문자로 변환
	for y := 0; y < newHeight; y++ {
		for x := 0; x < width; x++ {
			// 픽셀의 밝기 값(0-255)을 가져옵니다.
			pixelValue := gray.GetUCharAt(y, x)

			// 밝기 값을 charset 인덱스에 매핑합니다.
			charsetIndex := int(float64(pixelValue) / 255.0 * float64(len(c.charset)-1))
			buffer.WriteRune(c.charset[charsetIndex])
		}
		buffer.WriteRune('\n')
	}

	return buffer.String(), nil
}
