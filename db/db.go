package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

var globalDb *sql.DB

type queryResult struct {
	TableColumns []table.Column
	TableRows    []table.Row
}

func ConnectToDb(dbEngine, host, port, user, password, dbName string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch strings.ToLower(dbEngine) {
	case "postgres":
		db, err = connectToPostgres(host, port, user, password, dbName)
	case "mysql":
		db, err = connectToMysql(host, port, user, password, dbName)
	default:
		log.Fatal("You've entered an unsupported database engine. " +
			"The supported engines are postgres and mysql, " +
			"case insensitive.")

	}

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

	tableColumns := makeTableColumns(columns)
	tableRows := makeTableRows(rows, columns)

	return &queryResult{TableColumns: tableColumns, TableRows: tableRows}, nil
}
