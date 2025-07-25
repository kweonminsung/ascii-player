package tui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
	isPlaying    bool
	isPaused     bool
	frames       []string
	currentFrame int
}

// NewPlayer creates a new TUI player
func NewPlayer(filename string, fps int, loop bool, resolution string, color bool) *Player {
	app := tview.NewApplication()

	// Create text view for content
	textView := tview.NewTextView().
		SetDynamicColors(color).
		SetRegions(true).
		SetWordWrap(false).
		SetScrollable(false)
	textView.SetBorder(true).
		SetTitle(fmt.Sprintf(" ASCII Player - %s ", filename)).
		SetTitleAlign(tview.AlignLeft)

	// Create status view for controls
	statusView := tview.NewTextView().
		SetText("Controls: [SPACE] Pause/Resume | [R] Restart | [Q/ESC] Quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	statusView.SetBorder(true).SetTitle("Controls")

	return &Player{
		app:          app,
		textView:     textView,
		statusView:   statusView,
		fps:          fps,
		loop:         loop,
		resolution:   resolution,
		color:        color,
		filename:     filename,
		isPlaying:    false,
		isPaused:     false,
		frames:       []string{},
		currentFrame: 0,
	}
}

// LoadFrames loads frames for playback (placeholder for now)
func (p *Player) LoadFrames() error {
	// TODO: Implement actual frame loading from MP4
	// For now, create some dummy frames for demonstration
	p.frames = []string{
		"Frame 1: Loading " + p.filename + "...",
		"Frame 2: Processing video...",
		"Frame 3: Converting to ASCII...",
		"Frame 4: Playing animation...",
		"Frame 5: [ASCII art would be here]",
	}
	return nil
}

// Play starts the TUI player
func (p *Player) Play() error {
	if err := p.LoadFrames(); err != nil {
		return fmt.Errorf("failed to load frames: %v", err)
	}

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

	// Start playback loop
	go p.playbackLoop()

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
				go p.displayCurrentFrame() // 비동기적으로 호출
				return nil
			case 'r', 'R': // Restart
				p.currentFrame = 0
				p.isPaused = false
				// 이미 재생 중이면 새로운 goroutine을 생성하지 않음
				if !p.isPlaying {
					p.isPlaying = true
					go p.playbackLoop()
				}
				go p.displayCurrentFrame() // 비동기적으로 화면 업데이트
				return nil
			}
		}
		return event
	})
}

// playbackLoop handles the frame playback loop
func (p *Player) playbackLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(p.fps))
	defer ticker.Stop()

	p.isPlaying = true

	for range ticker.C {
		if !p.isPlaying {
			return
		}
		if !p.isPaused {
			p.displayCurrentFrame()
			p.currentFrame++

			if p.currentFrame >= len(p.frames) {
				if p.loop {
					p.currentFrame = 0
				} else {
					p.isPlaying = false
					p.displayEndMessage()
					return
				}
			}
		}
	}
}

// displayCurrentFrame shows the current frame
func (p *Player) displayCurrentFrame() {
	if p.currentFrame < len(p.frames) {
		status := "PLAYING"
		if p.isPaused {
			status = "PAUSED"
		}

		content := fmt.Sprintf("%s\n\n"+
			"═══════════════════════════════════════════════════════════\n"+
			"Frame: %d/%d | FPS: %d | Status: %s | Resolution: %s\n"+
			"═══════════════════════════════════════════════════════════",
			p.frames[p.currentFrame],
			p.currentFrame+1,
			len(p.frames),
			p.fps,
			status,
			p.resolution)

		p.app.QueueUpdateDraw(func() {
			p.textView.SetText(content)
		})
	}
}

// displayEndMessage shows the end message
func (p *Player) displayEndMessage() {
	content := fmt.Sprintf("╔════════════════════════════════════════════════════════════╗\n"+
		"║                     PLAYBACK FINISHED                     ║\n"+
		"╠════════════════════════════════════════════════════════════╣\n"+
		"║ File: %-49s║\n"+
		"║ Total frames: %-40d║\n"+
		"║                                                            ║\n"+
		"║ Press [R] to restart or [Q] to quit                       ║\n"+
		"╚════════════════════════════════════════════════════════════╝",
		p.filename, len(p.frames))

	p.app.QueueUpdateDraw(func() {
		p.textView.SetText(content)
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
