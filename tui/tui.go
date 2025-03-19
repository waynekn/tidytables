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
	return nil
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
		updatedModel, cmd := m.queryInput.Update(msg)
		m.queryInput = updatedModel.(*queryInput)
		return m, cmd
	} else if m.mode == viewTable {
		updatedModel, cmd := m.table.Update(msg)
		m.table = updatedModel.(*Table)
		return m, cmd
	} else if m.mode == viewError {
		updatedModel, cmd := m.errDisplay.Update(msg)
		m.errDisplay = updatedModel.(*errorDisplay)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() string {
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

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
