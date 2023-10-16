package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func initializeMySQL(database Database) (*sql.DB, error) {
	d, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&collation=utf8_unicode_ci&loc=%s&parseTime=true",
			database.Username,
			database.Password,
			database.Host,
			database.Port,
			database.DBName,
			url.QueryEscape(time.Local.String()),
		),
	)
	if err != nil {
		log.Panicln(err)
	}

	d.SetMaxOpenConns(database.MaxOpenConnections)
	d.SetMaxIdleConns(database.MaxIdleConnections)
	d.SetConnMaxLifetime(5 * time.Minute)

	if err := d.Ping(); err != nil {
		return nil, err
	}

	return d, nil
}
