package db

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/table"
	"github.com/waynekn/tidytables/utils"
)

var globalDb *sql.DB

type queryResult struct {
	TableColumns []table.Column
	TableRows    []table.Row
}

func ConnectToDb(host, port, user, password, dbName string) (*sql.DB, error) {
	db, err := ConnectPostgres(host, port, user, password, dbName)
	if err != nil {
		return nil, err
	}
	globalDb = db
	return db, nil

}

func QueryDB(q string) (*queryResult, error) {
	rows, err := globalDb.Query(q)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	tableColumns := utils.GetTableColumns(columns)
	tableRows := utils.GetTableRows(rows, columns)

	return &queryResult{TableColumns: tableColumns, TableRows: tableRows}, nil
}
