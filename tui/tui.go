package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
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

type queryResult struct {
	tableColumns []table.Column
	tableRows    []table.Row
}

type model struct {
	queryInput *queryInput
	table      *Table
	errDisplay *errorDisplay
	mode       viewMode
}

func initialModel() *model {
	return &model{
		queryInput: queryModel(),
		table:      initializeTable(),
		errDisplay: initErrorDisplay(),
		mode:       viewQuery,
	}
}

func (m *model) Init() tea.Cmd {
	// makes the textareas cursor blink
	return m.queryInput.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case QuerySubmittedMsg:
		res, err := db.QueryDB(msg.Query)

		if err != nil {
			m.mode = viewError
			return m, func() tea.Msg {
				return dbError{err: err}
			}

		}

		m.mode = viewTable
		return m, func() tea.Msg {
			return queryResult{
				tableColumns: res.TableColumns,
				tableRows:    res.TableRows,
			}
		}

	case switchToQueryInput:
		m.mode = viewQuery
		return m, m.queryInput.Init()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	if m.mode == viewQuery {
		_, cmd := m.queryInput.Update(msg)
		return m, cmd
	} else if m.mode == viewTable {
		_, cmd := m.table.Update(msg)
		return m, cmd
	} else if m.mode == viewError {
		_, cmd := m.errDisplay.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() string {
	switch m.mode {
	case viewQuery:
		return m.queryInput.View()
	case viewTable:
		return m.table.View()
	case viewError:
		return m.errDisplay.View()
	default:
		return ""
	}
}

func StartTea() {

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
