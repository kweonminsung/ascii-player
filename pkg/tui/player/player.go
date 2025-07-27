package player

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/ascii"
	"github.com/kweonminsung/ascii-player/pkg/audio"
	"github.com/kweonminsung/ascii-player/pkg/pixel"
	"github.com/kweonminsung/ascii-player/pkg/types"
	"github.com/kweonminsung/ascii-player/pkg/utils"
)

// getPlayerModeTitle returns the display title for the given mode
func getPlayerModeTitle(mode string) string {
	switch mode {
	case "pixel":
		return "PIXEL"
	case "ascii":
		return "ASCII"
	default:
		return "ASCII"
	}
}

// Player represents the TUI player
type Player struct {
	screen       tcell.Screen
	mutex        sync.Mutex
	fps          int
	loop         bool
	color        bool
	filename     string
	mode         string
	isPlaying    bool
	isPaused     bool
	isFinished   bool
	currentFrame int
	startTime    time.Time
	width        int
	height       int
	videoWidth   int
	videoHeight  int

	frameCountSinceLastCheck int
	lastFPSTime              time.Time
	actualFPS                float64

	// Player instances for different modes
	asciiPlayer *ascii.AsciiPlayer
	pixelPlayer *pixel.PixelPlayer
	audioPlayer *audio.AudioPlayer
}

// NewPlayer creates a new TUI player
func NewPlayer(filename string, fps int, loop bool, color bool, mode string) *Player {
	return &Player{
		fps:          fps,
		loop:         loop,
		color:        color,
		filename:     filename,
		mode:         mode,
		isPlaying:    false,
		isPaused:     false,
		isFinished:   false,
		currentFrame: 0,
	}
}

// GetFPS returns the FPS of the video.
func (p *Player) GetFPS() float64 {
	var fps float64
	switch p.mode {
	case "pixel":
		if p.pixelPlayer != nil {
			fps = p.pixelPlayer.GetFPS()
		}
	case "ascii":
		fallthrough
	default:
		if p.asciiPlayer != nil {
			fps = p.asciiPlayer.GetFPS()
		}
	}

	if fps > 0 {
		return fps
	}
	if p.fps > 0 {
		return float64(p.fps)
	}
	return 30 // Default FPS
}

// LoadFrames loads frames for playback
func (p *Player) LoadFrames() error {
	width, height := p.screen.Size()
	p.width, p.height = width, height-1 // Subtract 1 for status bar

	isYouTube := utils.IsValidYouTubeURL(p.filename)

	audioPlayer, err := audio.NewAudioPlayer(p.filename, isYouTube)
	if err != nil {
		log.Printf("failed to create audio player: %v. playing without audio", err)
	}
	p.audioPlayer = audioPlayer

	switch p.mode {
	case "pixel":
		pixelPlayer, err := pixel.NewPixelPlayer(p.filename, types.PlayerConfig{
			Mode:      "pixel",
			Color:     p.color,
			Width:     p.width,
			Height:    p.height,
			FPS:       p.fps,
			Loop:      p.loop,
			Source:    p.filename,
			IsYouTube: isYouTube,
		})
		if err != nil {
			return fmt.Errorf("failed to create pixel player: %v", err)
		}
		p.pixelPlayer = pixelPlayer
		p.videoWidth = pixelPlayer.GetVideoWidth()
		p.videoHeight = pixelPlayer.GetVideoHeight()
		return nil
	case "ascii":
		fallthrough
	default:
		asciiPlayer, err := ascii.NewAsciiPlayer(p.filename, types.PlayerConfig{
			Mode:      "ascii",
			Color:     p.color,
			Width:     p.width,
			Height:    p.height,
			FPS:       p.fps,
			Loop:      p.loop,
			Source:    p.filename,
			IsYouTube: isYouTube,
		})
		if err != nil {
			return fmt.Errorf("failed to create ASCII player: %v", err)
		}
		p.asciiPlayer = asciiPlayer
		p.videoWidth = asciiPlayer.GetVideoWidth()
		p.videoHeight = asciiPlayer.GetVideoHeight()
		return nil
	}
}

// Play starts the TUI player
func (p *Player) Play() error {
	var err error
	p.screen, err = tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to create screen: %v", err)
	}
	if err := p.screen.Init(); err != nil {
		return fmt.Errorf("failed to initialize screen: %v", err)
	}
	defer p.screen.Fini()
	if p.audioPlayer != nil {
		defer p.audioPlayer.Close()
	}

	if err := p.LoadFrames(); err != nil {
		return fmt.Errorf("failed to load frames: %v", err)
	}

	p.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	p.screen.Clear()

	// Handle interrupt signals
	go p.handleInterrupt()
	// Handle keyboard events
	go p.handleEvents()

	p.playbackLoop()

	return nil
}

func (p *Player) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.isPlaying = false
	if p.audioPlayer != nil {
		p.audioPlayer.Close()
	}
	p.screen.Fini()
	os.Exit(0)
}

func (p *Player) handleEvents() {
	for {
		ev := p.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			p.screen.Sync()
			width, height := p.screen.Size()
			p.width, p.height = width, height-1 // Subtract 1 for status bar

			switch p.mode {
			case "pixel":
				if p.pixelPlayer != nil {
					p.pixelPlayer.UpdateSize(p.width, p.height)
				}
			case "ascii":
				fallthrough
			default:
				if p.asciiPlayer != nil {
					p.asciiPlayer.UpdateSize(p.width, p.height)
				}
			}
			p.screen.Clear()
		case *tcell.EventKey:
			if p.isFinished {
				if ev.Rune() == 'r' || ev.Rune() == 'R' {
					p.startTime = time.Now()
					p.currentFrame = 0
					p.isPaused = false
					p.isFinished = false
					p.isPlaying = true
					p.screen.Clear()
				} else if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
					p.isPlaying = false
					p.screen.Fini()
					os.Exit(0)
				}
			} else {
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
					p.isPlaying = false
					p.screen.Fini()
					os.Exit(0)
				} else if ev.Rune() == ' ' {
					p.isPaused = !p.isPaused
					if p.audioPlayer != nil {
						if p.isPaused {
							p.audioPlayer.Pause()
						} else {
							p.audioPlayer.Resume()
						}
					}
				} else if ev.Rune() == 'r' || ev.Rune() == 'R' {
					p.startTime = time.Now()
					p.currentFrame = 0
					p.isPaused = false
					if p.audioPlayer != nil {
						p.audioPlayer.Rewind()
					}
				} else if ev.Key() == tcell.KeyRight {
					p.seek(5 * time.Second)
				} else if ev.Key() == tcell.KeyLeft {
					p.seek(-5 * time.Second)
				}
			}
		}
	}
}

func (p *Player) seek(duration time.Duration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch p.mode {
	case "pixel":
		if p.pixelPlayer != nil {
			p.pixelPlayer.Seek(duration)
		}
	case "ascii":
		fallthrough
	default:
		if p.asciiPlayer != nil {
			p.asciiPlayer.Seek(duration)
		}
	}
	if p.audioPlayer != nil {
		if err := p.audioPlayer.Seek(duration); err != nil {
			log.Printf("failed to seek audio: %v", err)
		}
	}
}

func (p *Player) playbackLoop() {
	var getNextFrame func() (string, error)
	var getFPS func() float64
	var getCurrentFrame func() int
	var getTotalFrames func() int

	switch p.mode {
	case "pixel":
		getNextFrame = p.pixelPlayer.GetNextFrame
		getFPS = p.pixelPlayer.GetFPS
		getCurrentFrame = p.pixelPlayer.GetCurrentFrame
		getTotalFrames = p.pixelPlayer.GetTotalFrames
	default:
		getNextFrame = p.asciiPlayer.GetNextFrame
		getFPS = p.asciiPlayer.GetFPS
		getCurrentFrame = p.asciiPlayer.GetCurrentFrame
		getTotalFrames = p.asciiPlayer.GetTotalFrames
	}

	fps := p.fps
	if fps <= 0 {
		fps = int(getFPS())
		if fps <= 0 {
			fps = 30
		}
	}

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	p.isPlaying = true
	p.startTime = time.Now()
	p.lastFPSTime = time.Now()

	if p.audioPlayer != nil {
		go p.audioPlayer.Play()
	}

	for range ticker.C {
		if !p.isPlaying {
			if p.isFinished {
				continue
			}
			return
		}
		if !p.isPaused {
			p.mutex.Lock()
			frame, err := getNextFrame()
			p.mutex.Unlock()
			if err != nil {
				if p.loop {
					// Reset by reloading frames
					if loadErr := p.LoadFrames(); loadErr != nil {
						log.Printf("failed to reload frames for looping: %v", loadErr)
						p.isPlaying = false
						return
					}
					p.startTime = time.Now()
					p.currentFrame = 0
					if p.audioPlayer != nil {
						p.audioPlayer.Rewind()
					}
					continue
				} else {
					p.isPlaying = false
					p.isFinished = true
					p.displayEndMessage()
					continue
				}
			}

			p.drawFrame(frame)
			p.currentFrame++
			p.frameCountSinceLastCheck++
		}

		if time.Since(p.lastFPSTime) >= time.Second {
			p.actualFPS = float64(p.frameCountSinceLastCheck) / time.Since(p.lastFPSTime).Seconds()
			p.lastFPSTime = time.Now()
			p.frameCountSinceLastCheck = 0
		}

		p.drawStatus(getCurrentFrame, getTotalFrames)
		p.screen.Show()
	}
}

func (p *Player) drawFrame(frame string) {
	p.screen.Clear()
	lines := strings.Split(frame, "\n")
	for y, line := range lines {
		p.drawString(0, y, line)
	}
}

func (p *Player) drawString(x, y int, str string) {
	style := tcell.StyleDefault
	runes := []rune(str)
	i := 0
	for i < len(runes) {
		r := runes[i]
		if r == '[' && i+8 < len(runes) && runes[i+1] == '#' && runes[i+8] == ']' {
			hex := string(runes[i+2 : i+8])
			if rgb, err := strconv.ParseInt(hex, 16, 32); err == nil {
				color := tcell.NewHexColor(int32(rgb))
				style = style.Foreground(color)
			}
			i += 9 // Move index past the color tag
			continue
		}
		p.screen.SetContent(x, y, r, nil, style)
		x++
		i++
	}
}

func (p *Player) drawStatus(getCurrentFrame func() int, getTotalFrames func() int) {
	_, screenHeight := p.screen.Size()
	statusY := p.height
	if statusY >= screenHeight {
		statusY = screenHeight - 1
	}

	status := "PLAYING"
	if p.isPaused {
		status = "PAUSED"
	}
	if !p.isPlaying {
		status = "STOPPED"
	}

	mode := "Normal"
	if p.loop {
		mode = "Loop"
	}

	currentFrame := getCurrentFrame()
	totalFrames := getTotalFrames()
	currentTime := time.Duration(float64(currentFrame)/p.GetFPS()) * time.Second
	totalTime := time.Duration(float64(totalFrames)/p.GetFPS()) * time.Second

	statusText := fmt.Sprintf("Mode: %s | FPS: %.1f/%d | Status: %s | Frame: %d/%d | Time: %s/%s | Resolution: %s | Player: %s | Controls: [SPACE] Pause/Resume | [R] Restart | [<-/->] Seek | [Q/ESC] Quit",
		mode,
		p.actualFPS,
		p.fps,
		status,
		currentFrame,
		totalFrames,
		utils.FormatDuration(currentTime),
		utils.FormatDuration(totalTime),
		strconv.Itoa(p.width)+"x"+strconv.Itoa(p.height),
		getPlayerModeTitle(p.mode))

	// Clear status line
	width, _ := p.screen.Size()
	for i := 0; i < width; i++ {
		p.screen.SetContent(i, statusY, ' ', nil, tcell.StyleDefault.Background(tcell.ColorSilver).Foreground(tcell.ColorBlack))
	}

	// Draw status text
	runes := []rune(statusText)
	for i := 0; i < len(runes) && i < width; i++ {
		p.screen.SetContent(i, statusY, runes[i], nil, tcell.StyleDefault.Background(tcell.ColorSilver).Foreground(tcell.ColorBlack))
	}
}

func (p *Player) displayEndMessage() {
	p.screen.Clear()
	screenWidth, screenHeight := p.screen.Size()

	boxWidth := 50
	boxHeight := 7

	x := (screenWidth - boxWidth) / 2
	y := (screenHeight - boxHeight) / 2

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	// Draw the box
	for i := x; i < x+boxWidth; i++ {
		p.screen.SetContent(i, y, tcell.RuneHLine, nil, style)
		p.screen.SetContent(i, y+boxHeight-1, tcell.RuneHLine, nil, style)
	}
	for i := y; i < y+boxHeight; i++ {
		p.screen.SetContent(x, i, tcell.RuneVLine, nil, style)
		p.screen.SetContent(x+boxWidth-1, i, tcell.RuneVLine, nil, style)
	}
	p.screen.SetContent(x, y, tcell.RuneULCorner, nil, style)
	p.screen.SetContent(x+boxWidth-1, y, tcell.RuneURCorner, nil, style)
	p.screen.SetContent(x, y+boxHeight-1, tcell.RuneLLCorner, nil, style)
	p.screen.SetContent(x+boxWidth-1, y+boxHeight-1, tcell.RuneLRCorner, nil, style)

	// Fill the box
	for i := y + 1; i < y+boxHeight-1; i++ {
		for j := x + 1; j < x+boxWidth-1; j++ {
			p.screen.SetContent(j, i, ' ', nil, style)
		}
	}

	// Draw the message
	msg1 := "Playback Finished!"
	msg2 := "Do you want to restart?"
	msg3 := "[R] Restart  [Q] Quit"
	p.drawString(x+(boxWidth-len(msg1))/2, y+2, msg1)
	p.drawString(x+(boxWidth-len(msg2))/2, y+3, msg2)
	p.drawString(x+(boxWidth-len(msg3))/2, y+5, msg3)

	p.screen.Show()
}
