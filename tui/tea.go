package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/waynekn/tidytables/db"
)

type viewMode int

const (
	viewQuery viewMode = iota
	viewTable
)

type model struct {
	queryInput queryInput
	table      Table
	mode       viewMode
}

func initialModel() model {
	return model{
		queryInput: queryModel(),
		table:      initializeTable(),
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
			log.Fatal(color.RedString(err.Error()))
		}
		m.table = NewTable(res.TableColumns, res.TableRows)
		m.mode = viewTable
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlT:
			if m.mode == viewQuery {
				m.mode = viewTable
			} else {
				m.mode = viewQuery
			}

			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	if m.mode == viewQuery {
		var cmd tea.Cmd
		var qi tea.Model
		qi, cmd = m.queryInput.Update(msg)
		m.queryInput = qi.(queryInput)
		return m, cmd

	} else if m.mode == viewTable {
		var cmd tea.Cmd
		var t tea.Model
		t, cmd = m.table.Update(msg)
		m.table = t.(Table)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.mode == viewQuery {
		return m.queryInput.View()
	} else if m.mode == viewTable {
		return m.table.View()
	}
	return ""
}

func StartTea() {

	p := tea.NewProgram(initialModel() /*tea.WithAltScreen()*/)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
