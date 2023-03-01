package config

import (
	"database/sql"
	"log"
)

func initializeSQLLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./erp.db")
	if err != nil {
		log.Panicln(err)
	}

	return db, nil
}

// // insert
// stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
// log.Panicln(err)
