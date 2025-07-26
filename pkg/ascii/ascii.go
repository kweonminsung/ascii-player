package ascii

import (
	"fmt"
	"log"
	"time"

	"github.com/kweonminsung/ascii-player/pkg/media"
	"github.com/kweonminsung/ascii-player/pkg/types"
)

// Player represents an ASCII video player
type AsciiPlayer struct {
	extractor *media.FrameExtractor
	converter *media.AsciiConverter
	config    types.PlayerConfig
}

// NewAsciiPlayer creates a new ASCII player instance
func NewAsciiPlayer(source string, config types.PlayerConfig) (*AsciiPlayer, error) {
	extractor, err := media.NewFrameExtractor(source, config.IsYouTube)
	if err != nil {
		return nil, fmt.Errorf("failed to create frame extractor: %v", err)
	}

	converter := media.NewAsciiConverter()

	return &AsciiPlayer{
		extractor: extractor,
		converter: converter,
		config:    config,
	}, nil
}

// Close closes the ASCII player and releases resources
func (p *AsciiPlayer) Close() {
	if p.extractor != nil {
		p.extractor.Close()
	}
}

// GetFPS returns the FPS of the video
func (p *AsciiPlayer) GetFPS() float64 {
	return p.extractor.GetFPS()
}

// GetFrameAt seeks to a specific time and returns the ASCII art for that frame
func (p *AsciiPlayer) GetFrameAt(seekTime time.Duration) (string, error) {
	frame, err := p.extractor.GetFrameAt(seekTime)
	if err != nil {
		return "", fmt.Errorf("failed to get frame at %v: %v", seekTime, err)
	}
	defer frame.Close()

	if frame.Empty() {
		return "", fmt.Errorf("got empty frame at %v", seekTime)
	}

	asciiArt, err := p.converter.Convert(frame, p.config.Width, p.config.Height, p.config.Color)
	if err != nil {
		return "", fmt.Errorf("failed to convert frame to ASCII: %v", err)
	}

	return asciiArt, nil
}

// PlayConsecutiveFrames plays consecutive frames as ASCII art in the console
func (p *AsciiPlayer) PlayConsecutiveFrames(frameCount int) error {
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

		asciiArt, err := p.converter.Convert(frame, p.config.Width, p.config.Height, p.config.Color)
		frame.Close()

		if err != nil {
			log.Printf("Failed to convert frame %d: %v", i, err)
			continue
		}

		// Clear terminal and print ASCII art
		fmt.Printf("\033[2J\033[H%s", asciiArt)
		time.Sleep(frameInterval)
	}

	return nil
}

// PlayFrameAtTime seeks to a specific time and displays the ASCII art frame
func (p *AsciiPlayer) PlayFrameAtTime(seekTime time.Duration) error {
	asciiArt, err := p.GetFrameAt(seekTime)
	if err != nil {
		return err
	}

	// Clear terminal and print ASCII art
	fmt.Printf("\033[2J\033[H%s", asciiArt)
	return nil
}

// GetVideoInfo returns basic information about the video
func (p *AsciiPlayer) GetVideoInfo() map[string]interface{} {
	return map[string]interface{}{
		"fps":    p.GetFPS(),
		"width":  p.config.Width,
		"height": p.config.Height,
	}
}

// PlayYouTubeVideo is a convenience function to play a YouTube video as ASCII art
func PlayYouTubeVideo(youtubeURL string, width, height int, seekTime time.Duration, frameCount int) error {
	player, err := NewAsciiPlayer(youtubeURL,
		types.PlayerConfig{
			Width:     width,
			Height:    height,
			FPS:       30,
			Loop:      false,
			Source:    youtubeURL,
			IsYouTube: true,
		})
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	defer player.Close()

	log.Printf("Successfully opened YouTube video. FPS: %.2f", player.GetFPS())

	// Seek to specific time if provided
	if seekTime > 0 {
		log.Printf("Seeking to %v...", seekTime)
		err := player.PlayFrameAtTime(seekTime)
		if err != nil {
			return fmt.Errorf("failed to play frame at %v: %v", seekTime, err)
		}
		time.Sleep(2 * time.Second) // Display for 2 seconds
	}

	// Play consecutive frames
	if frameCount > 0 {
		log.Printf("Playing %d consecutive frames...", frameCount)
		err := player.PlayConsecutiveFrames(frameCount)
		if err != nil {
			return fmt.Errorf("failed to play consecutive frames: %v", err)
		}
	}

	return nil
}

// PlayLocalVideo is a convenience function to play a local video file as ASCII art
func PlayLocalVideo(filePath string, width, height int, frameCount int) error {
	player, err := NewAsciiPlayer(filePath, types.PlayerConfig{
		Width:     width,
		Height:    height,
		FPS:       30,
		Loop:      false,
		Source:    filePath,
		IsYouTube: false,
	})
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	defer player.Close()

	log.Printf("Successfully opened local video. FPS: %.2f", player.GetFPS())

	// Play consecutive frames
	err = player.PlayConsecutiveFrames(frameCount)
	if err != nil {
		return fmt.Errorf("failed to play consecutive frames: %v", err)
	}

	return nil
}
