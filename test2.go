package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/media"
	"github.com/kweonminsung/ascii-player/pkg/pixel"
	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run testv2.go <video_file_or_youtube_url>")
		return
	}
	source := os.Args[1]
	fps := 30
	loop := false
	isYouTube := isValidYouTubeURL(source)

	playPixel(source, isYouTube, fps, loop)
}

func isValidYouTubeURL(url string) bool {
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

func playPixel(source string, isYouTube bool, fps int, loop bool) {
	extractor, err := media.NewFrameExtractor(source, isYouTube)
	if err != nil {
		log.Fatalf("Error creating frame extractor: %v", err)
	}
	defer extractor.Close()

	if fps <= 0 {
		fps = int(extractor.GetFPS())
	}
	if fps <= 0 {
		fps = 30 // Default FPS
	}

	converter, err := pixel.NewColorConverter()
	if err != nil {
		log.Fatalf("Error creating ANSI converter: %v", err)
	}
	defer converter.Close()

	// Handle quit event
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
			break
		}
		if frame.Empty() {
			frame.Close()
			continue
		}

		startTime := time.Now()

		termWidth, termHeight := converter.GetScreenSize()
		resizedFrame := gocv.NewMat()
		// Resize frame to fit terminal, preserving aspect ratio
		img, _ := frame.ToImage()
		aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
		newWidth := termWidth / 3 // Adjust for 3 characters per pixel
		newHeight := int(float64(newWidth) / aspectRatio)
		if newHeight > termHeight {
			newHeight = termHeight
			newWidth = int(float64(newHeight) * aspectRatio)
		}

		gocv.Resize(frame, &resizedFrame, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationDefault)
		frame.Close()

		if err := converter.DrawFrame(resizedFrame); err != nil {
			log.Printf("Error drawing frame: %v", err)
		}
		resizedFrame.Close()

		elapsed := time.Since(startTime)
		if elapsed < frameDuration {
			time.Sleep(frameDuration - elapsed)
		}
	}
}
