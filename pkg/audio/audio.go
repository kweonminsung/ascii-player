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
	ytdlp "github.com/kweonminsung/console-cinema/third_party/yt-dlp"
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
	speed     float64
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
		speed:     1.0,
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
	executablePath, err := ytdlp.GetExecutablePath()
	if err != nil {
		return err
	}
	defer os.Remove(executablePath)

	cmd := exec.Command(
		executablePath,
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
	ap.speed = speed
	if ap.resampler != nil {
		speaker.Lock()
		ap.resampler.SetRatio(ap.speed)
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
	speaker.Lock()
	defer speaker.Unlock()

	if err := ap.streamer.Seek(0); err != nil {
		return err
	}

	// Re-create the resampler to clear its internal state.
	ap.resampler = beep.Resample(4, ap.format.SampleRate, ap.format.SampleRate, ap.streamer)
	ap.resampler.SetRatio(ap.speed)
	ap.ctrl.Streamer = ap.resampler
	return nil
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
	if ap.streamer.Len() > 0 && newPosition >= ap.streamer.Len() {
		newPosition = ap.streamer.Len() - 1
	} else if ap.streamer.Len() <= 0 {
		newPosition = 0
	}

	if err := ap.streamer.Seek(newPosition); err != nil {
		return err
	}

	// Re-create the resampler to clear its internal state after a seek.
	ap.resampler = beep.Resample(4, ap.format.SampleRate, ap.format.SampleRate, ap.streamer)
	ap.resampler.SetRatio(ap.speed)
	ap.ctrl.Streamer = ap.resampler

	return nil
}

// SeekAbsolute seeks the audio to an absolute position specified by the duration.
func (ap *AudioPlayer) SeekAbsolute(duration time.Duration) error {
	speaker.Lock()
	defer speaker.Unlock()

	if duration < 0 {
		duration = 0
	}

	newPosition := ap.format.SampleRate.N(duration)
	if ap.streamer.Len() > 0 && newPosition >= ap.streamer.Len() {
		newPosition = ap.streamer.Len() - 1
	}

	if err := ap.streamer.Seek(newPosition); err != nil {
		return err
	}

	// Re-create the resampler to clear its internal state after a seek.
	ap.resampler = beep.Resample(4, ap.format.SampleRate, ap.format.SampleRate, ap.streamer)
	ap.resampler.SetRatio(ap.speed)
	ap.ctrl.Streamer = ap.resampler

	return nil
}

// Close closes the audio player and cleans up resources
func (ap *AudioPlayer) Close() {
	if ap.closer != nil {
		ap.closer.Close()
	}
	speaker.Close()
	// The audio file is now managed by the cache, so we don't remove it here.
}
