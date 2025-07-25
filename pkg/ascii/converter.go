package ascii

import (
	"bytes"
	"image"

	"gocv.io/x/gocv"
)

// DefaultCharset은 ASCII 변환에 사용될 기본 문자 집합입니다.
// 밝기 순서에 따라 문자를 정렬합니다.
const DefaultCharset = " .:-=+*#%@"

// Converter는 이미지를 ASCII 아트로 변환하는 기능을 제공합니다.
type Converter struct {
	charset []rune
}

// NewConverter는 새로운 Converter 인스턴스를 생성합니다.
// charset이 비어있으면 DefaultCharset을 사용합니다.
func NewConverter(charset string) *Converter {
	if charset == "" {
		charset = DefaultCharset
	}
	return &Converter{
		charset: []rune(charset),
	}
}

// Convert는 gocv.Mat 이미지를 지정된 너비와 높이의 ASCII 문자열로 변환합니다.
func (c *Converter) Convert(img gocv.Mat, width, height int) (string, error) {
	// 최종 ASCII 아트 출력을 위한 버퍼
	var buffer bytes.Buffer

	// 1. 이미지 크기 조절
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: width, Y: height}, 0, 0, gocv.InterpolationLinear)

	// 2. 흑백으로 변환
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)

	// 3. 픽셀을 순회하며 ASCII 문자로 변환
	for y := 0; y < height; y++ {
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
