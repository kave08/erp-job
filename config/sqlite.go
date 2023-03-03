package config

import (
	"database/sql"
	"log"
)

func initializeSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./erp.db")
	if err != nil {
		log.Panicln(err)
	}

	return db, nil
}