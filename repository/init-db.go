package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	connStr := "user=postgres password=p@ssw0rd dbname=mydatabase sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
