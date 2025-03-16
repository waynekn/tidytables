package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type dbError struct {
	err error
}

type errorDisplay struct {
	textarea textarea.Model
	err      error
}

// returns the width of the terminal window - 5
func getDisplayWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	// substract -5 for some margin
	return width - 5
}

func initErrorDisplay() errorDisplay {
	ti := textarea.New()

	ti.FocusedStyle = textarea.Style{
		Text:       lipgloss.NewStyle().Foreground(lipgloss.Color("#f50202")),
		CursorLine: lipgloss.NewStyle().Foreground(lipgloss.Color("#fc1703")),
		Prompt:     lipgloss.NewStyle().Foreground(lipgloss.Color("#fc1703")),
	}

	ti.SetWidth(getDisplayWidth())

	ti.Focus()

	return errorDisplay{
		textarea: ti,
		err:      nil,
	}
}

func (m errorDisplay) Init() tea.Cmd {
	return nil
}

func (m errorDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case dbError:
		m.textarea.SetValue(msg.err.Error())
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlBackslash:
			return m, func() tea.Msg {
				return switchToQueryInput{}
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m errorDisplay) View() string {

	return fmt.Sprintf(
		"%s\n\n%s",
		m.textarea.View(),
		"(type ctrl+\\ to edit query)",
	) + "\n\n"
}
