package main

import (
	"log"
	"time"

	"github.com/kweonminsung/ascii-player/pkg/ascii"
)

func main() {
	// --- YouTube 스트리밍 처리 예제 (간단한 버전) ---
	youtubeURL := "https://www.youtube.com/watch?v=OZytLLasceA" // Rick Astley - Never Gonna Give You Up
	log.Printf("--- Processing YouTube stream: %s ---\n", youtubeURL)

	// 간편한 방법: 내장 함수 사용
	err := ascii.PlayYouTubeVideo(
		youtubeURL,
		120,            // width
		40,             // height
		43*time.Second, // seek to 43 seconds
		300,            // play 300 frames (about 10 seconds)
	)
	if err != nil {
		log.Fatalf("Failed to play YouTube video: %v", err)
	}

	log.Println("\nTask completed successfully.")
}

// --- 고급 사용법 예제 (주석 처리) ---
/*
func advancedExample() {
	youtubeURL := "https://www.youtube.com/watch?v=OZytLLasceA"

	// 1. ASCII Player 생성
	player, err := ascii.NewPlayer(youtubeURL, true, 120, 40)
	if err != nil {
		log.Fatalf("Failed to create player: %v", err)
	}
	defer player.Close()

	log.Printf("Successfully opened video. FPS: %.2f\n", player.GetFPS())

	// 2. 특정 시간의 프레임 가져오기
	seekTime := 43 * time.Second
	log.Printf("Seeking to %v...", seekTime)

	err = player.PlayFrameAtTime(seekTime)
	if err != nil {
		log.Fatalf("Failed to play frame at %v: %v", seekTime, err)
	}

	time.Sleep(2 * time.Second) // 2초간 표시

	// 3. 연속 프레임 재생
	log.Println("Playing consecutive frames...")
	err = player.PlayConsecutiveFrames(300)
	if err != nil {
		log.Printf("Error during playback: %v", err)
	}
}

func localVideoExample() {
	// 로컬 비디오 파일 재생
	err := ascii.PlayLocalVideo("video.mp4", 120, 40, 300)
	if err != nil {
		log.Fatalf("Failed to play local video: %v", err)
	}
}
*/
