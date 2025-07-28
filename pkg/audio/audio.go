package audio

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// AudioPlayer manages audio playback
type AudioPlayer struct {
	ctrl      *beep.Ctrl
	streamer  beep.StreamSeeker
	format    beep.Format
	closer    io.Closer
	audioPath string
	resampler *beep.Resampler
	mutex     sync.Mutex
}

// NewAudioPlayer creates a new AudioPlayer
func NewAudioPlayer(videoPath string, isYouTube bool) (*AudioPlayer, error) {
	audioPath := fmt.Sprintf("temp_audio_%d.mp3", time.Now().UnixNano())
	var err error

	if isYouTube {
		err = extractAudioFromYouTube(videoPath, audioPath)
		if err != nil {
			return nil, fmt.Errorf("failed to extract audio from YouTube: %v", err)
		}
	} else {
		err = extractAudio(videoPath, audioPath)
		if err != nil {
			return nil, fmt.Errorf("failed to extract audio: %v", err)
		}
	}

	f, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to decode mp3: %v", err)
	}

	resampler := beep.Resample(4, format.SampleRate, format.SampleRate, streamer)
	ctrl := &beep.Ctrl{Streamer: resampler, Paused: false}

	return &AudioPlayer{
		ctrl:      ctrl,
		streamer:  streamer,
		format:    format,
		closer:    f,
		audioPath: audioPath,
		resampler: resampler,
	}, nil
}

func extractAudio(videoPath, audioPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-q:a", "0", "-map", "a", audioPath, "-y")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func extractAudioFromYouTube(videoURL, audioPath string) error {
	cmd := exec.Command(
		"yt-dlp",
		"--no-playlist",
		"--quiet",
		"--progress",
		"-x",
		"--audio-format", "mp3",
		"-o", audioPath,
		videoURL,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("yt-dlp command failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// Play starts audio playback
func (ap *AudioPlayer) Play() {
	speaker.Init(ap.format.SampleRate, ap.format.SampleRate.N(time.Second/10))
	speaker.Play(ap.ctrl)
}

// SetSpeed adjusts the playback speed of the audio.
func (ap *AudioPlayer) SetSpeed(speed float64) {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()
	if ap.resampler != nil {
		speaker.Lock()
		ap.resampler.SetRatio(speed)
		speaker.Unlock()
	}
}

// Pause pauses audio playback
func (ap *AudioPlayer) Pause() {
	speaker.Lock()
	ap.ctrl.Paused = true
	speaker.Unlock()
}

// Resume resumes audio playback
func (ap *AudioPlayer) Resume() {
	speaker.Lock()
	ap.ctrl.Paused = false
	speaker.Unlock()
}

// Rewind rewinds the audio to the beginning
func (ap *AudioPlayer) Rewind() error {
	return ap.streamer.Seek(0)
}

// Seek seeks the audio by the given duration.
func (ap *AudioPlayer) Seek(duration time.Duration) error {
	speaker.Lock()
	defer speaker.Unlock()

	currentPosition := ap.streamer.Position()
	currentDuration := ap.format.SampleRate.D(currentPosition)
	newDuration := currentDuration + duration

	if newDuration < 0 {
		newDuration = 0
	}

	newPosition := ap.format.SampleRate.N(newDuration)
	if newPosition >= ap.streamer.Len() {
		newPosition = ap.streamer.Len() - 1
	}

	return ap.streamer.Seek(newPosition)
}

// Close closes the audio player and cleans up resources
func (ap *AudioPlayer) Close() {
	if ap.closer != nil {
		ap.closer.Close()
	}
	speaker.Close()
	if ap.audioPath != "" {
		os.Remove(ap.audioPath)
	}
}
