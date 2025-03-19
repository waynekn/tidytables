package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type (
	errMsg error
)

type Styles struct {
	BorderColor lipgloss.Color
}

type QuerySubmittedMsg struct {
	Query string
}

func defaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#0390fc")
	return s
}

type queryInput struct {
	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	styles      *Styles
}

func queryModel() *queryInput {
	ta := textarea.New()
	ta.Placeholder = "Write a query"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to Tidy Tables!
Easily view SQL query results in a clean, visually appealing table right from your terminal.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &queryInput{
		textarea:    ta,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		styles:      defaultStyles(),
	}
}

func (m *queryInput) Init() tea.Cmd {
	return textarea.Blink
}

func (m *queryInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		m.viewport.GotoBottom()

		// subtract 2 columns to leave some margin on each side
		newWidth := msg.Width - 2
		m.textarea.SetWidth(newWidth)
		m.viewport.Width = newWidth
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:

			return m, tea.Quit
		case tea.KeyCtrlN:
			m.textarea.InsertString("\n")
		case tea.KeyEnter:
			value := m.textarea.Value()
			value = strings.TrimSpace(value)

			if value == "" {
				return m, nil
			}

			q := strings.ReplaceAll(value, "\n", " ")
			q = strings.ReplaceAll(q, "\"", "'")

			return m, func() tea.Msg {
				return QuerySubmittedMsg{Query: q}
			}

		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *queryInput) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.styles.BorderColor).
		MarginRight(5)

	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		borderStyle.Render(m.textarea.View()),
	)
}
