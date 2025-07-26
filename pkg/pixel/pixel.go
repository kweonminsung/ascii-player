package pixel

import (
	"fmt"
	"image"
	"log"
	"regexp"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/media"
	"gocv.io/x/gocv"
)

// PixelPlayer represents a pixel-based video player
type PixelPlayer struct {
	extractor *media.FrameExtractor
	converter *ColorConverter
	fps       int
	loop      bool
	source    string
	isYouTube bool
}

// PlayerConfig holds configuration for the pixel player
type PlayerConfig struct {
	FPS    int
	Loop   bool
	Source string
}

// NewPixelPlayer creates a new pixel player instance
func NewPixelPlayer(source string, config PlayerConfig) (*PixelPlayer, error) {
	isYouTube := IsValidYouTubeURL(source)

	extractor, err := media.NewFrameExtractor(source, isYouTube)
	if err != nil {
		return nil, fmt.Errorf("error creating frame extractor: %v", err)
	}

	fps := config.FPS
	if fps <= 0 {
		fps = int(extractor.GetFPS())
	}
	if fps <= 0 {
		fps = 30 // Default FPS
	}

	converter, err := NewColorConverter()
	if err != nil {
		extractor.Close()
		return nil, fmt.Errorf("error creating ANSI converter: %v", err)
	}

	return &PixelPlayer{
		extractor: extractor,
		converter: converter,
		fps:       fps,
		loop:      config.Loop,
		source:    source,
		isYouTube: isYouTube,
	}, nil
}

// Close closes the pixel player and releases resources
func (p *PixelPlayer) Close() {
	if p.converter != nil {
		p.converter.Close()
	}
	if p.extractor != nil {
		p.extractor.Close()
	}
}

// GetFPS returns the current FPS setting
func (p *PixelPlayer) GetFPS() int {
	return p.fps
}

// SetFPS sets the FPS for playback
func (p *PixelPlayer) SetFPS(fps int) {
	if fps > 0 {
		p.fps = fps
	}
}

// GetVideoFPS returns the original video FPS
func (p *PixelPlayer) GetVideoFPS() float64 {
	if p.extractor != nil {
		return p.extractor.GetFPS()
	}
	return 0
}

// IsLooping returns whether the player is set to loop
func (p *PixelPlayer) IsLooping() bool {
	return p.loop
}

// SetLoop sets the loop mode
func (p *PixelPlayer) SetLoop(loop bool) {
	p.loop = loop
}

// GetScreenSize returns the current screen dimensions
func (p *PixelPlayer) GetScreenSize() (int, int) {
	if p.converter != nil {
		return p.converter.GetScreenSize()
	}
	return 0, 0
}

// Play starts playing the video
func (p *PixelPlayer) Play() error {
	return p.PlayWithCallback(nil)
}

// PlayWithCallback plays the video with a custom frame callback
func (p *PixelPlayer) PlayWithCallback(frameCallback func()) error {
	// Handle quit events
	go p.handleInput()

	frameDuration := time.Second / time.Duration(p.fps)

	for {
		frame, err := p.extractor.ReadNextFrame()
		if err != nil {
			if p.loop {
				if seekErr := p.extractor.Seek(0); seekErr != nil {
					log.Printf("Failed to seek to beginning: %v", seekErr)
					break
				}
				continue
			}
			break
		}

		if frame.Empty() {
			frame.Close()
			continue
		}

		startTime := time.Now()

		if err := p.processAndDrawFrame(frame); err != nil {
			log.Printf("Error processing frame: %v", err)
		}
		frame.Close()

		// Call custom callback if provided
		if frameCallback != nil {
			frameCallback()
		}

		elapsed := time.Since(startTime)
		if elapsed < frameDuration {
			time.Sleep(frameDuration - elapsed)
		}
	}

	return nil
}

// processAndDrawFrame processes a frame and draws it to the screen
func (p *PixelPlayer) processAndDrawFrame(frame gocv.Mat) error {
	termWidth, termHeight := p.converter.GetScreenSize()
	resizedFrame := gocv.NewMat()
	defer resizedFrame.Close()

	// Resize frame to fit terminal, preserving aspect ratio
	img, err := frame.ToImage()
	if err != nil {
		return fmt.Errorf("failed to convert frame to image: %v", err)
	}

	aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	newWidth := termWidth / 3 // Adjust for 3 characters per pixel
	newHeight := int(float64(newWidth) / aspectRatio)
	if newHeight > termHeight {
		newHeight = termHeight
		newWidth = int(float64(newHeight) * aspectRatio)
	}

	gocv.Resize(frame, &resizedFrame, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationDefault)

	return p.converter.DrawFrame(resizedFrame)
}

// handleInput handles keyboard input events
func (p *PixelPlayer) handleInput() {
	for {
		if p.converter == nil || p.converter.Screen() == nil {
			break
		}

		ev := p.converter.Screen().PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				p.Close()
				return
			}
		}
	}
}

// Seek seeks to a specific position in the video
func (p *PixelPlayer) Seek(position time.Duration) error {
	if p.extractor != nil {
		return p.extractor.Seek(position)
	}
	return fmt.Errorf("no extractor available")
}

// GetPlayerInfo returns information about the player
func (p *PixelPlayer) GetPlayerInfo() map[string]interface{} {
	termWidth, termHeight := p.GetScreenSize()
	return map[string]interface{}{
		"source":       p.source,
		"isYouTube":    p.isYouTube,
		"fps":          p.fps,
		"videoFPS":     p.GetVideoFPS(),
		"loop":         p.loop,
		"screenWidth":  termWidth,
		"screenHeight": termHeight,
	}
}

// IsValidYouTubeURL checks if the given URL is a valid YouTube URL
func IsValidYouTubeURL(url string) bool {
	patterns := []string{
		`^https?://(www\.)?youtube\.com/watch\?v=[\w-]+`,
		`^https?://youtu\.be/[\w-]+`,
		`^https?://(www\.)?youtube\.com/embed/[\w-]+`,
	}
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}
	return false
}

// PlayVideo is a convenience function to play a video with simple parameters
func PlayVideo(source string, fps int, loop bool) error {
	config := PlayerConfig{
		FPS:  fps,
		Loop: loop,
	}

	player, err := NewPixelPlayer(source, config)
	if err != nil {
		return fmt.Errorf("failed to create pixel player: %v", err)
	}
	defer player.Close()

	log.Printf("Playing video: %s (FPS: %d, Loop: %t)", source, player.GetFPS(), player.IsLooping())
	if player.isYouTube {
		log.Printf("Detected YouTube video")
	}

	return player.Play()
}

// PlayVideoWithCallback is a convenience function to play a video with a custom callback
func PlayVideoWithCallback(source string, fps int, loop bool, callback func()) error {
	config := PlayerConfig{
		FPS:  fps,
		Loop: loop,
	}

	player, err := NewPixelPlayer(source, config)
	if err != nil {
		return fmt.Errorf("failed to create pixel player: %v", err)
	}
	defer player.Close()

	log.Printf("Playing video: %s (FPS: %d, Loop: %t)", source, player.GetFPS(), player.IsLooping())
	return player.PlayWithCallback(callback)
}
