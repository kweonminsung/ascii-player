package player

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/ascii-player/pkg/ascii"
	"github.com/kweonminsung/ascii-player/pkg/pixel"
	"github.com/kweonminsung/ascii-player/pkg/types"
	"github.com/kweonminsung/ascii-player/pkg/utils"
	"github.com/rivo/tview"
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
	app          *tview.Application
	textView     *tview.TextView
	statusView   *tview.TextView
	fps          int
	loop         bool
	resolution   string
	color        bool
	filename     string
	mode         string
	isPlaying    bool
	isPaused     bool
	frames       []string
	currentFrame int
	// Player instances for different modes
	asciiPlayer *ascii.AsciiPlayer
	pixelPlayer *pixel.PixelPlayer
}

// NewPlayer creates a new TUI player
func NewPlayer(filename string, fps int, loop bool, resolution string, color bool, mode string) *Player {
	app := tview.NewApplication()

	// Create text view for content
	textView := tview.NewTextView().
		SetDynamicColors(color).
		SetRegions(true).
		SetWordWrap(false).
		SetScrollable(false)
	textView.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s Player - %s ", getPlayerModeTitle(mode), filename)).
		SetTitleAlign(tview.AlignLeft)

	// Create status view for controls
	statusView := tview.NewTextView().
		SetText("Loading...").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	statusView.SetBorder(true).SetTitle("Status")

	return &Player{
		app:          app,
		textView:     textView,
		statusView:   statusView,
		fps:          fps,
		loop:         loop,
		resolution:   resolution,
		color:        color,
		filename:     filename,
		mode:         mode,
		isPlaying:    false,
		isPaused:     false,
		frames:       []string{},
		currentFrame: 0,
	}
}

// LoadFrames loads frames for playback
func (p *Player) LoadFrames() error {
	width, height := 120, 40 // Default values
	if p.resolution != "" {
		fmt.Sscanf(p.resolution, "%dx%d", &width, &height)
	}

	// Check if it's a YouTube URL
	isYouTube := utils.IsValidYouTubeURL(p.filename)

	switch p.mode {
	case "pixel":
		pixelPlayer, err := pixel.NewPixelPlayer(p.filename, types.PlayerConfig{
			Mode:      "pixel",
			Color:     p.color,
			Width:     width,
			Height:    height,
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
			Width:     width,
			Height:    height,
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
	if err := p.LoadFrames(); err != nil {
		return fmt.Errorf("failed to load frames: %v", err)
	}

	// 모든 모드에서 TUI를 사용
	return p.playWithTUI()
}

// playWithTUI plays video using TUI interface for ASCII mode
func (p *Player) playWithTUI() error {
	// Set up key bindings
	p.setupKeyBindings()

	// Set up the main layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.textView, 0, 1, true).
		AddItem(p.statusView, 3, 0, false)

	p.app.SetRoot(flex, true)

	// Handle interrupt signals
	go p.handleInterrupt()

	// Start playback loop for 해당 모드
	if p.mode == "pixel" {
		go p.pixelPlaybackLoop()
	} else {
		go p.asciiPlaybackLoop()
	}

	// Initialize status view
	go p.updateStatusView()

	// Run the application
	return p.app.Run()
}

// setupKeyBindings sets up keyboard controls
func (p *Player) setupKeyBindings() {
	p.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			p.app.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				p.app.Stop()
				return nil
			case ' ': // Space bar to pause/resume
				p.isPaused = !p.isPaused
				go p.updateStatusView() // 상태 업데이트
				return nil
			case 'r', 'R': // Restart
				if p.mode == "ascii" {
					// ASCII 모드 재시작
					if p.isPlaying || p.isPaused {
						p.currentFrame = 0
						p.isPaused = false
						if !p.isPlaying {
							p.isPlaying = true
							go p.asciiPlaybackLoop()
						}
						go p.updateStatusView()
					}
				} else if p.mode == "pixel" {
					// Pixel 모드 재시작
					if p.isPlaying || p.isPaused {
						p.currentFrame = 0
						p.isPaused = false
						if !p.isPlaying {
							p.isPlaying = true
							go p.pixelPlaybackLoop()
						}
						go p.updateStatusView()
					}
				}
				return nil
			}
		}
		return event
	})
}

// updateStatusView updates the status view with current information
func (p *Player) updateStatusView() {
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

	frameInfo := ""
	if p.mode == "ascii" || p.mode == "pixel" {
		frameInfo = fmt.Sprintf("Frame: %d/%d", p.currentFrame+1, len(p.frames))
	} else {
		frameInfo = "Frame: N/A"
	}

	statusText := fmt.Sprintf("Mode: %s | %s | FPS: %d | Status: %s | Resolution: %s | Player: %s | Controls: [SPACE] Pause/Resume | [R] Restart | [Q/ESC] Quit",
		mode,
		frameInfo,
		p.fps,
		status,
		p.resolution,
		getPlayerModeTitle(p.mode))

	p.app.QueueUpdateDraw(func() {
		p.statusView.SetText(statusText)
	})
}

// asciiPlaybackLoop handles frame playback for ASCII mode
func (p *Player) asciiPlaybackLoop() {
	if p.asciiPlayer == nil {
		return
	}

	fps := p.fps
	if fps <= 0 {
		fps = int(p.asciiPlayer.GetFPS())
		if fps <= 0 {
			fps = 30
		}
	}

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	p.isPlaying = true
	currentTime := time.Duration(0)
	frameDuration := time.Second / time.Duration(fps)

	for range ticker.C {
		if !p.isPlaying {
			return
		}
		if !p.isPaused {
			frame, err := p.asciiPlayer.GetFrameAt(currentTime)
			if err != nil {
				log.Printf("Error getting frame: %v", err)
				continue
			}

			p.app.QueueUpdateDraw(func() {
				p.textView.SetText(frame)
			})

			p.currentFrame++
			p.updateStatusView()
			currentTime += frameDuration

			// Check if we've reached the end (this is a simple check, might need improvement)
			if currentTime > time.Minute*10 { // Arbitrary limit, should be replaced with actual video duration
				if p.loop {
					currentTime = 0
					p.currentFrame = 0
				} else {
					p.isPlaying = false
					p.displayEndMessage()
					p.updateStatusView()
					return
				}
			}
		}
	}
}

// pixelPlaybackLoop handles frame playback for Pixel mode
func (p *Player) pixelPlaybackLoop() {
	if p.pixelPlayer == nil {
		return
	}

	fps := p.fps
	if fps <= 0 {
		fps = int(p.pixelPlayer.GetFPS())
		if fps <= 0 {
			fps = 30
		}
	}

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	p.isPlaying = true
	currentTime := time.Duration(0)
	frameDuration := time.Second / time.Duration(fps)

	for range ticker.C {
		if !p.isPlaying {
			return
		}
		if !p.isPaused {
			frame, err := p.pixelPlayer.GetFrameAt(currentTime)
			if err != nil {
				log.Printf("Error getting pixel frame: %v", err)
				continue
			}

			p.app.QueueUpdateDraw(func() {
				p.textView.SetText(frame)
			})

			p.currentFrame++
			p.updateStatusView()
			currentTime += frameDuration

			// Check if we've reached the end
			if currentTime > time.Minute*10 { // Arbitrary limit, should be replaced with actual video duration
				if p.loop {
					currentTime = 0
					p.currentFrame = 0
				} else {
					p.isPlaying = false
					p.displayEndMessage()
					p.updateStatusView()
					return
				}
			}
		}
	}
}

// displayEndMessage shows the end message
func (p *Player) displayEndMessage() {
	// Create a modal for end message
	totalFrames := "N/A"
	if p.mode == "ascii" {
		totalFrames = fmt.Sprintf("%d", len(p.frames))
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Playback Finished!\n\nFile: %s\nTotal frames: %s\n\nPress [R] to restart or [Q] to quit",
			p.filename, totalFrames)).
		AddButtons([]string{"Restart (R)", "Quit (Q)"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 || buttonLabel == "Restart (R)" {
				// Restart
				p.currentFrame = 0
				p.isPaused = false
				p.isPlaying = true
				p.app.SetRoot(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(p.textView, 0, 1, true).
					AddItem(p.statusView, 3, 0, false), true)
				if p.mode == "ascii" {
					go p.asciiPlaybackLoop()
				}
				go p.updateStatusView()
			} else {
				// Quit
				p.app.Stop()
			}
		})

	p.app.QueueUpdateDraw(func() {
		p.app.SetRoot(modal, true)
	})
}

// handleInterrupt handles interrupt signals
func (p *Player) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.isPlaying = false
	p.app.Stop()
}
