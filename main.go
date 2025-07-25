package main

import (
	"log"
	"time"

	"github.com/kweonminsung/ascii-player/pkg/ascii"
	"github.com/kweonminsung/ascii-player/pkg/media"
)

func main() {
	// --- YouTube 스트리밍 처리 예제 ---
	youtubeURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ" // Rick Astley - Never Gonna Give You Up
	log.Printf("--- Processing YouTube stream: %s ---\n", youtubeURL)

	// 1. FrameExtractor 생성
	extractor, err := media.NewFrameExtractor(youtubeURL, true)
	if err != nil {
		log.Fatalf("Failed to create frame extractor: %v", err)
	}
	defer extractor.Close()

	log.Printf("Successfully opened video. FPS: %.2f\n", extractor.GetFPS())

	// 2. 특정 시간(예: 43초, "Never gonna give you up" 부분)으로 이동하여 프레임 가져오기
	seekTime := 43 * time.Second
	log.Printf("Seeking to %v...", seekTime)

	frame, err := extractor.GetFrameAt(seekTime)
	if err != nil {
		log.Fatalf("Failed to get frame at %v: %v", seekTime, err)
	}

	// ASCII 변환기 생성
	converter := ascii.NewConverter("")

	if !frame.Empty() {
		log.Printf("Successfully got frame at %v. Frame size: %s\n", seekTime, frame.Size())
		// 프레임을 ASCII로 변환
		asciiArt, err := converter.Convert(frame, 120, 40) // 터미널 크기에 맞게 너비, 높이 조절
		if err != nil {
			log.Fatalf("Failed to convert frame to ASCII: %v", err)
		}
		// 터미널을 지우고 ASCII 아트 출력
		log.Printf("\033[2J\033[H%s", asciiArt)
	}
	frame.Close() // 프레임 사용 후 반드시 닫기

	// 3. 1초씩 앞으로 가면서 5개 프레임 연속으로 읽고 ASCII로 출력
	log.Println("\nPlaying 5 consecutive frames as ASCII art...")
	frameInterval := time.Second / time.Duration(extractor.GetFPS())

	for i := 0; i < 300; i++ { // 300프레임 (약 10초) 재생
		loopFrame, err := extractor.ReadNextFrame()
		if err != nil {
			log.Printf("Could not read next frame: %v", err)
			break
		}
		if loopFrame.Empty() {
			log.Println("Got an empty frame, end of stream?")
			loopFrame.Close()
			break
		}

		asciiArt, err := converter.Convert(loopFrame, 120, 40)
		if err != nil {
			log.Printf("Failed to convert frame: %v", err)
			loopFrame.Close()
			continue
		}

		// 터미널을 지우고 ASCII 아트 출력
		log.Printf("\033[2J\033[H%s", asciiArt)
		loopFrame.Close() // 각 프레임 사용 후 바로 닫기

		time.Sleep(frameInterval) // FPS에 맞춰 잠시 대기
	}

	log.Println("\nTask completed successfully.")
}
