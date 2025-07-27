package pixel

import (
	"fmt"
	"log"
	"time"

	"github.com/kweonminsung/ascii-player/pkg/media"
	"github.com/kweonminsung/ascii-player/pkg/types"
)

// PixelPlayer represents a pixel-based video player
type PixelPlayer struct {
	extractor *media.FrameExtractor
	converter *media.AnsiConverter
	config    types.PlayerConfig
}

// NewPixelPlayer creates a new pixel player instance
func NewPixelPlayer(source string, config types.PlayerConfig) (*PixelPlayer, error) {
	extractor, err := media.NewFrameExtractor(source, config.IsYouTube)
	if err != nil {
		return nil, fmt.Errorf("failed to create frame extractor: %v", err)
	}

	converter := media.NewAnsiConverter()

	return &PixelPlayer{
		extractor: extractor,
		converter: converter,
		config:    config,
	}, nil
}

// Close closes the player and releases resources
func (p *PixelPlayer) Close() {
	if p.extractor != nil {
		p.extractor.Close()
	}
}

// GetFPS returns the FPS of the video
func (p *PixelPlayer) GetFPS() float64 {
	return p.extractor.GetFPS()
}

// GetFrameAt seeks to a specific time and returns the ANSI-colored ASCII art for that frame
func (p *PixelPlayer) GetFrameAt(seekTime time.Duration) (string, error) {
	frame, err := p.extractor.GetFrameAt(seekTime)
	if err != nil {
		return "", fmt.Errorf("failed to get frame at %v: %v", seekTime, err)
	}
	defer frame.Close()

	if frame.Empty() {
		return "", fmt.Errorf("got empty frame at %v", seekTime)
	}

	pixelArt, err := p.converter.Convert(frame, p.config.Width, p.config.Height, p.config.Color)
	if err != nil {
		return "", fmt.Errorf("failed to convert frame to pixel art: %v", err)
	}

	return pixelArt, nil
}

// PlayConsecutiveFrames plays consecutive frames as ANSI-colored pixel art in the console
func (p *PixelPlayer) PlayConsecutiveFrames(frameCount int) error {
	frameInterval := time.Second / time.Duration(p.GetFPS())

	for i := 0; i < frameCount; i++ {
		frame, err := p.extractor.ReadNextFrame()
		if err != nil {
			return fmt.Errorf("could not read frame %d: %v", i, err)
		}

		if frame.Empty() {
			frame.Close()
			log.Println("Got an empty frame, end of stream")
			break
		}

		pixelArt, err := p.converter.Convert(frame, p.config.Width, p.config.Height, p.config.Color)
		frame.Close()

		if err != nil {
			log.Printf("Failed to convert frame %d: %v", i, err)
			continue
		}

		// Clear terminal and print pixel art
		fmt.Printf("\033[2J\033[H%s", pixelArt)
		time.Sleep(frameInterval)
	}

	return nil
}

// PlayFrameAtTime seeks to a specific time and displays the ANSI-colored frame
func (p *PixelPlayer) PlayFrameAtTime(seekTime time.Duration) error {
	pixelArt, err := p.GetFrameAt(seekTime)
	if err != nil {
		return err
	}

	fmt.Printf("\033[2J\033[H%s", pixelArt)
	return nil
}

// GetVideoInfo returns basic information about the video
func (p *PixelPlayer) GetVideoInfo() map[string]interface{} {
	return map[string]interface{}{
		"fps":    p.GetFPS(),
		"width":  p.config.Width,
		"height": p.config.Height,
	}
}

// GetVideoWidth returns the original width of the video.
func (p *PixelPlayer) GetVideoWidth() int {
	return p.extractor.GetWidth()
}

// GetVideoHeight returns the original height of the video.
func (p *PixelPlayer) GetVideoHeight() int {
	return p.extractor.GetHeight()
}

// UpdateSize updates the player's dimensions.
func (p *PixelPlayer) UpdateSize(width, height int) {
	p.config.Width = width
	p.config.Height = height
}

// GetNextFrame reads the next frame and converts it to pixel art.
func (p *PixelPlayer) GetNextFrame() (string, error) {
	frame, err := p.extractor.ReadNextFrame()
	if err != nil {
		return "", fmt.Errorf("could not read next frame: %v", err)
	}
	defer frame.Close()

	if frame.Empty() {
		return "", fmt.Errorf("got empty frame")
	}

	pixelArt, err := p.converter.Convert(frame, p.config.Width, p.config.Height, p.config.Color)
	if err != nil {
		return "", fmt.Errorf("failed to convert frame to pixel art: %v", err)
	}

	return pixelArt, nil
}

// Seek seeks the video by the given duration.
func (p *PixelPlayer) Seek(duration time.Duration) {
	currentPos := p.extractor.GetPosition()
	newPos := currentPos + duration
	if newPos < 0 {
		newPos = 0
	}
	p.extractor.Seek(newPos)
}

// GetCurrentFrame returns the current frame number of the video.
func (p *PixelPlayer) GetCurrentFrame() int {
	return p.extractor.GetCurrentFrame()
}

// GetTotalFrames returns the total number of frames in the video.
func (p *PixelPlayer) GetTotalFrames() int {
	return p.extractor.GetTotalFrames()
}
