package types

// YScaleFactor는 터미널 문자의 세로/가로 비율을 보정하기 위한 값입니다.
// 일반적인 터미널 폰트는 높이가 너비의 약 2배이므로 0.5에 가까운 값을 사용합니다.
const YScaleFactor = 0.55

// PlayerConfig holds configuration
type PlayerConfig struct {
	Mode      string
	Color     bool
	FPS       int
	Width     int
	Height    int
	Loop      bool
	Source    string
	IsYouTube bool
}
