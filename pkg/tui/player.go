package tui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Player represents the TUI player
type Player struct {
	tui          *TUI
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
	tui, err := NewTUI()
	if err != nil {
		return nil // In case of error, return nil (should handle this better)
	}

	return &Player{
		tui:          tui,
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
	defer p.tui.Close()

	if err := p.LoadFrames(); err != nil {
		return fmt.Errorf("failed to load frames: %v", err)
	}

	// Handle interrupt signals
	go p.handleInterrupt()

	// Start playback loop
	go p.playbackLoop()

	// Handle keyboard input
	p.handleInput()

	return nil
}

// setupKeyBindings and createControlsView are no longer needed with basic TUI

// playFrames handles the frame playback loop
func (p *Player) playbackLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(p.fps))
	defer ticker.Stop()

	p.isPlaying = true

	for p.isPlaying {
		select {
		case <-ticker.C:
			if !p.isPaused {
				p.displayCurrentFrame()
				p.currentFrame++

				if p.currentFrame >= len(p.frames) {
					if p.loop {
						p.currentFrame = 0
					} else {
						p.isPlaying = false
						p.displayEndMessage()
					}
				}
			}
		default:
			if !p.isPlaying {
				return
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
			"Controls: [SPACE] Pause/Resume | [R] Restart | [Q] Quit\n"+
			"═══════════════════════════════════════════════════════════",
			p.frames[p.currentFrame],
			p.currentFrame+1,
			len(p.frames),
			p.fps,
			status,
			p.resolution)
		p.tui.Display(content)
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
	p.tui.Display(content)
}

// showCurrentFrame displays the current frame
func (p *Player) showCurrentFrame() {
	p.displayCurrentFrame()
}

// handleInput handles keyboard input
func (p *Player) handleInput() {
	for {
		switch ev := p.tui.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				p.isPlaying = false
				return
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q', 'Q':
					p.isPlaying = false
					return
				case ' ': // Space bar for pause/resume
					p.isPaused = !p.isPaused
					p.displayCurrentFrame()
				case 'r', 'R': // Restart
					p.currentFrame = 0
					p.isPaused = false
					if !p.isPlaying {
						p.isPlaying = true
						go p.playbackLoop()
					}
				}
			}
		case *tcell.EventResize:
			p.tui.screen.Sync()
			p.displayCurrentFrame()
		}
	}
} // handleInterrupt handles interrupt signals
func (p *Player) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.isPlaying = false
	p.tui.Close()
	os.Exit(0)
}
