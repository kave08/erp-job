package repository

import (
	"database/sql"
	"erp-job/repository/database"
)

type Repository struct {
	Database database.DatabaseInterface
}

func NewRepository(sdb *sql.DB) *Repository {
	return &Repository{
		Database: database.NewDatabase(sdb),
	}
}
