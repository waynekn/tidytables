package utils

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func GetTableRows(rows *sql.Rows, columns []string) []table.Row {
	// Create a slice to hold values for each row
	var tableRows []table.Row

	// Prepare a slice for holding the values during scanning
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err)
		}

		// Format date
		for i := range values {
			if reflect.TypeOf(values[i]) == reflect.TypeOf(time.Time{}) {
				t := values[i].(time.Time)
				// 2 January, 2006 03:04 PM
				values[i] = t.Format("2/1/2006")
			}
		}

		// slice of strings to represent the row
		row := make([]string, len(columns))
		for i, value := range values {
			row[i] = fmt.Sprintf("%v", value)
		}

		tableRows = append(tableRows, row)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tableRows
}
