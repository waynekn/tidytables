package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waynekn/tidytables/db"
)

type viewMode int

const (
	viewQuery viewMode = iota
	viewTable
	viewError
)

type switchToQueryInput struct{}

type model struct {
	queryInput queryInput
	table      Table
	errDisplay errorDisplay
	mode       viewMode
}

func initialModel() model {
	return model{
		queryInput: queryModel(),
		table:      initializeTable(),
		errDisplay: initErrorDisplay(),
		mode:       viewQuery,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case QuerySubmittedMsg:
		res, err := db.QueryDB(msg.Query)

		if err != nil {
			m.mode = viewError
			return m, func() tea.Msg {
				return dbError{err: err}
			}

		}

		m.table = NewTable(res.TableColumns, res.TableRows)
		m.mode = viewTable
	case switchToQueryInput:
		m.mode = viewQuery
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	if m.mode == viewQuery {
		var qi tea.Model
		qi, cmd = m.queryInput.Update(msg)
		m.queryInput = qi.(queryInput)
		return m, cmd

	} else if m.mode == viewTable {
		var t tea.Model
		t, cmd = m.table.Update(msg)
		m.table = t.(Table)
		return m, cmd
	} else if m.mode == viewError {
		var e tea.Model
		e, cmd = m.errDisplay.Update(msg)
		m.errDisplay = e.(errorDisplay)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.mode == viewQuery {
		return m.queryInput.View()
	} else if m.mode == viewTable {
		return m.table.View()
	} else if m.mode == viewError {
		return m.errDisplay.View()
	}
	return ""
}

func StartTea() {

	p := tea.NewProgram(initialModel() /*tea.WithAltScreen()*/)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
