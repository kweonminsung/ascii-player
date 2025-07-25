package media

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

// FrameExtractor는 비디오 소스를 열고 프레임을 제어하는 기능을 제공합니다.
type FrameExtractor struct {
	vc      *gocv.VideoCapture
	source  string
	isYouTube bool
	fps     float64
}

// NewFrameExtractor는 새로운 FrameExtractor 인스턴스를 생성합니다.
// isYouTube 플래그가 true이면, yt-dlp를 이용해 스트림 URL을 얻어옵니다.
func NewFrameExtractor(source string, isYouTube bool) (*FrameExtractor, error) {
	videoSource := source
	if isYouTube {
		// yt-dlp를 사용하여 가장 좋은 품질의 mp4 스트림 URL을 가져옵니다. (max 720p)
		cmd := exec.Command("yt-dlp", "-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4][height<=720]/best", "-g", source)
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute yt-dlp for url %s: %w. Is yt-dlp installed and in your PATH?", source, err)
		}

		streamURL := strings.TrimSpace(string(output))
		if streamURL == "" {
			return nil, fmt.Errorf("yt-dlp returned an empty stream URL")
		}
		// yt-dlp는 비디오와 오디오 URL을 별도 라인으로 반환할 수 있습니다. 첫 번째(비디오)를 사용합니다.
		videoSource = strings.Split(streamURL, "\n")[0]
	}

	vc, err := gocv.VideoCaptureFile(videoSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open video capture from source %s: %w", videoSource, err)
	}
	if !vc.IsOpened() {
		vc.Close()
		return nil, fmt.Errorf("video capture is not opened for source: %s", videoSource)
	}

	fps := vc.Get(gocv.VideoCaptureFPS)

	return &FrameExtractor{
		vc:      vc,
		source:  source,
		isYouTube: isYouTube,
		fps:     fps,
	}, nil
}

// ReadNextFrame은 현재 위치에서 다음 프레임을 읽어 반환합니다.
func (c *FrameExtractor) ReadNextFrame() (gocv.Mat, error) {
	mat := gocv.NewMat()
	if ok := c.vc.Read(&mat); !ok {
		mat.Close()
		return gocv.Mat{}, fmt.Errorf("failed to read frame or end of stream")
	}
	if mat.Empty() {
		mat.Close()
		return gocv.Mat{}, fmt.Errorf("read an empty frame")
	}
	return mat, nil
}

// Seek는 비디오의 재생 위치를 지정된 시간으로 이동시킵니다.
func (c *FrameExtractor) Seek(d time.Duration) error {
	ms := float64(d.Milliseconds())
	c.vc.Set(gocv.VideoCapturePosMsec, ms)
	return nil
}

// GetFrameAt은 지정된 시간의 프레임을 정확히 가져옵니다.
// 내부적으로 Seek 후 ReadNextFrame을 호출합니다.
func (c *FrameExtractor) GetFrameAt(d time.Duration) (gocv.Mat, error) {
	if err := c.Seek(d); err != nil {
		return gocv.Mat{}, err
	}
	c.vc.Grab(1) 
	return c.ReadNextFrame()
}

// GetFPS는 비디오의 FPS를 반환합니다.
func (c *FrameExtractor) GetFPS() float64 {
	return c.fps
}

// Close는 사용된 모든 리소스를 해제합니다.
func (c *FrameExtractor) Close() {
	if c.vc != nil {
		c.vc.Close()
	}
}
