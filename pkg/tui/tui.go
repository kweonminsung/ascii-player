package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// TUI represents the terminal user interface
type TUI struct {
	screen tcell.Screen
}

// NewTUI initializes a new TUI instance
func NewTUI() (*TUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	return &TUI{screen: screen}, nil
}

// Close cleans up the TUI resources
func (t *TUI) Close() {
	t.screen.Fini()
}

// Display renders the TUI content
func (t *TUI) Display(content string) {
	t.screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	width, height := t.screen.Size()
	lines := splitIntoLines(content, width, height)
	for i, line := range lines {
		for j, r := range []rune(line) {
			t.screen.SetContent(j, i, r, nil, style)
		}
	}
	t.screen.Show()
}

// splitIntoLines splits content into lines that fit the screen
func splitIntoLines(content string, width int, height int) []string {
	lines := strings.Split(content, "\n")
	result := []string{}

	for _, line := range lines {
		if len(line) <= width {
			result = append(result, line)
		} else {
			// Wrap long lines
			for len(line) > width {
				result = append(result, line[:width])
				line = line[width:]
			}
			if len(line) > 0 {
				result = append(result, line)
			}
		}

		if len(result) >= height {
			break
		}
	}

	return result
}
