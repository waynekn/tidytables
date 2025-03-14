package utils

import "github.com/charmbracelet/bubbles/table"

// Converts a slice of strings to a slice of table.Column
func GetTableColumns(cols []string) []table.Column {
	var tableColumns = make([]table.Column, 0)
	for _, col := range cols {
		tableColumns = append(tableColumns, table.Column{Title: col, Width: 10})
	}
	return tableColumns
}
