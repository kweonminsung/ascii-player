package player

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/ascii"
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
	fps          int
	loop         bool
	resolution   string
	color        bool
	filename     string
	mode         string
	isPlaying    bool
	isPaused     bool
	currentFrame int
	startTime    time.Time
	width        int
	height       int

	// Player instances for different modes
	asciiPlayer *ascii.AsciiPlayer
	pixelPlayer *pixel.PixelPlayer
}

// NewPlayer creates a new TUI player
func NewPlayer(filename string, fps int, loop bool, resolution string, color bool, mode string) *Player {
	return &Player{
		fps:          fps,
		loop:         loop,
		resolution:   resolution,
		color:        color,
		filename:     filename,
		mode:         mode,
		isPlaying:    false,
		isPaused:     false,
		currentFrame: 0,
	}
}

// LoadFrames loads frames for playback
func (p *Player) LoadFrames() error {
	p.width, p.height = 120, 40 // Default values
	if p.resolution != "" {
		fmt.Sscanf(p.resolution, "%dx%d", &p.width, &p.height)
	}

	isYouTube := utils.IsValidYouTubeURL(p.filename)

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

	// If no resolution is specified, use the full screen size
	if p.resolution == "" {
		width, height := p.screen.Size()
		p.resolution = fmt.Sprintf("%dx%d", width, height-1) // Subtract 1 for status bar
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
			p.resolution = fmt.Sprintf("%dx%d", width, height-1)
			if err := p.LoadFrames(); err != nil {
				log.Printf("failed to reload frames on resize: %v", err)
			}
			p.screen.Clear()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
				p.isPlaying = false
				p.screen.Fini()
				os.Exit(0)
			} else if ev.Rune() == ' ' {
				p.isPaused = !p.isPaused
			} else if ev.Rune() == 'r' || ev.Rune() == 'R' {
				p.startTime = time.Now()
				p.currentFrame = 0
				p.isPaused = false
			}
		}
	}
}

func (p *Player) playbackLoop() {
	var getFrame func(time.Duration) (string, error)
	var getFPS func() float64

	switch p.mode {
	case "pixel":
		getFrame = p.pixelPlayer.GetFrameAt
		getFPS = p.pixelPlayer.GetFPS
	default:
		getFrame = p.asciiPlayer.GetFrameAt
		getFPS = p.asciiPlayer.GetFPS
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

	for range ticker.C {
		if !p.isPlaying {
			return
		}
		if !p.isPaused {
			elapsed := time.Since(p.startTime)
			frame, err := getFrame(elapsed)
			if err != nil {
				if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "out of range") {
					if p.loop {
						p.startTime = time.Now()
						p.currentFrame = 0
						continue
					} else {
						p.isPlaying = false
						p.displayEndMessage()
						return
					}
				}
				log.Printf("Error getting frame: %v", err)
				continue
			}

			p.drawFrame(frame)
			p.drawStatus()
			p.screen.Show()

			p.currentFrame++
		}
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

func (p *Player) drawStatus() {
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

	statusText := fmt.Sprintf("Mode: %s | FPS: %d | Status: %s | Resolution: %s | Player: %s | Controls: [SPACE] Pause/Resume | [R] Restart | [Q/ESC] Quit",
		mode,
		p.fps,
		status,
		p.resolution,
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
	msg := "Playback Finished! Press [R] to restart or [Q] to quit."
	p.drawString(0, 0, msg)
	p.screen.Show()
}
