package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const defaultMySqlPort = "3306"

// Handles connection to a MySQL database.
//
// Parameters:
//   - host: the database host
//   - port: the database port
//   - user: the database user
//   - password: the database user's password
//   - dbName: the name of the database
//
// Returns:
//   - *sql.DB: a pointer to the database connection
//   - error: an error if the connection fails
func connectToMysql(host, port, user, password, dbName string) (*sql.DB, error) {

	var mySqlPort string

	if port == "" {
		mySqlPort = defaultMySqlPort
	} else {
		mySqlPort = port
	}

	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", host, mySqlPort),
		DBName: dbName,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
