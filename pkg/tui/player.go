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
				go p.displayCurrentFrame() // 비동기적으로 호출
				go p.updateStatusView()    // 상태 업데이트
				return nil
			case 'r', 'R': // Restart
				// 재생이 끝난 상태가 아닐 때만 재시작 처리
				if p.isPlaying || p.isPaused {
					p.currentFrame = 0
					p.isPaused = false
					// 이미 재생 중이면 새로운 goroutine을 생성하지 않음
					if !p.isPlaying {
						p.isPlaying = true
						go p.playbackLoop()
					}
					go p.displayCurrentFrame() // 비동기적으로 화면 업데이트
					go p.updateStatusView()    // 상태 업데이트
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

	statusText := fmt.Sprintf("Mode: %s | Frame: %d/%d | FPS: %d | Status: %s | Resolution: %s | Controls: [SPACE] Pause/Resume | [R] Restart | [Q/ESC] Quit",
		mode,
		p.currentFrame+1,
		len(p.frames),
		p.fps,
		status,
		p.resolution)

	p.app.QueueUpdateDraw(func() {
		p.statusView.SetText(statusText)
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
			p.updateStatusView() // 프레임 변경 시 상태 업데이트
			p.currentFrame++

			if p.currentFrame >= len(p.frames) {
				if p.loop {
					p.currentFrame = 0
				} else {
					p.isPlaying = false
					p.displayEndMessage()
					p.updateStatusView() // 종료 시 상태 업데이트
					return
				}
			}
		}
	}
}

// displayCurrentFrame shows the current frame
func (p *Player) displayCurrentFrame() {
	if p.currentFrame < len(p.frames) {
		content := fmt.Sprintf("%s\n\n%s",
			p.frames[p.currentFrame],
			"═══════════════════════════════════════════════════════════")

		p.app.QueueUpdateDraw(func() {
			p.textView.SetText(content)
		})
	}
}

// displayEndMessage shows the end message
func (p *Player) displayEndMessage() {
	// Create a modal for end message
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Playback Finished!\n\nFile: %s\nTotal frames: %d\n\nPress [R] to restart or [Q] to quit",
			p.filename, len(p.frames))).
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
				go p.playbackLoop()
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
