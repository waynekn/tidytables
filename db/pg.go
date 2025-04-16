package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultPgPort = "5432"

// Handles connection to a PostgreSQL database.
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
func connectToPostgres(host, port, user, password, dbName string) (*sql.DB, error) {
	var pgPort string

	if port == "" {
		pgPort = defaultPgPort
	} else {
		pgPort = port
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, pgPort, dbName)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, err
}
