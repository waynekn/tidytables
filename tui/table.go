package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Table struct {
	table table.Model
}

func (m *Table) Init() tea.Cmd {
	return nil
}

func (m *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case queryResult:

		newTable := createTable(msg.tableRows, msg.tableColumns)

		newTable.SetStyles(getTableStyles())
		m.table = newTable
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlBackslash:
			return m, func() tea.Msg {
				return switchToQueryInput{}
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Table) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		baseStyle.Render(m.table.View()),
		"(type ctrl+\\ to edit query)",
	)

}

func initializeTable() *Table {
	columns := []table.Column{}

	rows := []table.Row{}

	t := createTable(rows, columns)
	t.SetStyles(getTableStyles())

	return &Table{table: t}
}

func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}

func createTable(rows []table.Row, columns []table.Column) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	return t
}
