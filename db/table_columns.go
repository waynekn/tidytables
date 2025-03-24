package db

import "github.com/charmbracelet/bubbles/table"

// Converts a slice of column names (strings) into a slice of table.Column.
//
// Each column is initialized with a default width of 10.
func makeTableColumns(cols []string) []table.Column {
	var tableColumns = make([]table.Column, 0)
	for _, col := range cols {
		tableColumns = append(tableColumns, table.Column{Title: col, Width: 10})
	}
	return tableColumns
}
