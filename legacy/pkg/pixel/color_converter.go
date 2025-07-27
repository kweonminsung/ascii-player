package pixel

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/media"
	"gocv.io/x/gocv"
)

const YScaleFactor = 0.55

// ColorConverter는 gocv.Mat를 tcell 화면에 그리는 역할을 합니다.
type ColorConverter struct {
	screen tcell.Screen
}

// NewColorConverter는 새로운 ColorConverter를 생성합니다.
func NewColorConverter() (*ColorConverter, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	return &ColorConverter{screen: screen}, nil
}

// Close는 tcell 화면을 닫습니다.
func (c *ColorConverter) Close() {
	c.screen.Fini()
}

// DrawFrame은 gocv.Mat를 tcell 화면에 그립니다.
func (c *ColorConverter) DrawFrame(frame gocv.Mat) error {
	_ = image.Point{} // Force compiler to recognize image package usage
	img, err := frame.ToImage()
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Clear screen before drawing
	c.screen.Clear()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			color := tcell.NewRGBColor(int32(r>>8), int32(g>>8), int32(b>>8))
			style := tcell.StyleDefault.Background(color)
			// Use 3 characters for each pixel to adjust aspect ratio
			c.screen.SetContent(x*3, y, ' ', nil, style)
			c.screen.SetContent(x*3+1, y, ' ', nil, style)
			c.screen.SetContent(x*3+2, y, ' ', nil, style)
		}
	}

	c.screen.Show()
	return nil
}

// GetScreenSize는 tcell 화면의 크기를 반환합니다.
func (c *ColorConverter) GetScreenSize() (int, int) {
	return c.screen.Size()
}

// Screen은 tcell.Screen 객체를 반환합니다.
func (c *ColorConverter) Screen() tcell.Screen {
	return c.screen
}

// Play는 주어진 소스(파일 또는 유튜브 URL)를 ANSI 아트로 재생합니다.
func Play(source string, isYouTube bool, fps int, loop bool) {
	extractor, err := media.NewFrameExtractor(source, isYouTube)
	if err != nil {
		log.Fatalf("Error creating frame extractor: %v", err)
	}
	defer extractor.Close()

	if fps <= 0 {
		fps = int(extractor.GetFPS())
	}
	if fps <= 0 {
		fps = 30 // 기본 FPS
	}

	converter, err := NewColorConverter()
	if err != nil {
		log.Fatalf("Error creating ANSI converter: %v", err)
	}
	defer converter.Close()

	// 'q' 또는 Ctrl+C를 누르면 종료되도록 이벤트 핸들러 설정
	go func() {
		for {
			ev := converter.Screen().PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					converter.Close()
					os.Exit(0)
				}
			}
		}
	}()

	frameDuration := time.Second / time.Duration(fps)

	for {
		frame, err := extractor.ReadNextFrame()
		if err != nil {
			if loop {
				if seekErr := extractor.Seek(0); seekErr != nil {
					log.Printf("Failed to seek to beginning: %v", seekErr)
					break
				}
				continue
			}
			break // 루프가 아니면 종료
		}
		if frame.Empty() {
			frame.Close()
			continue
		}

		startTime := time.Now()

		termWidth, termHeight := converter.GetScreenSize()
		resizedFrame := gocv.NewMat()

		// 터미널 크기에 맞춰 프레임 크기 조절 (가로세로 비율 유지)
		img, _ := frame.ToImage()
		aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()) // width / height
		newWidth := termWidth / 3                                             // ANSI는 픽셀당 3문자를 사용하므로
		newHeight := int(float64(newWidth) / aspectRatio * YScaleFactor)

		// 높이가 터미널을 초과하면 높이를 기준으로 너비를 다시 계산
		if newHeight > termHeight {
			newHeight = termHeight
			newWidth = int(float64(newHeight) * aspectRatio / YScaleFactor)
		}

		gocv.Resize(frame, &resizedFrame, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationDefault)
		frame.Close() // 원본 프레임은 이제 필요 없으므로 닫기

		if err := converter.DrawFrame(resizedFrame); err != nil {
			log.Printf("Error drawing frame: %v", err)
		}
		resizedFrame.Close()

		elapsed := time.Since(startTime)
		if elapsed < frameDuration {
			time.Sleep(frameDuration - elapsed)
		}
	}
	fmt.Println("Playback finished.")
}
