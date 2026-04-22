package app

import (
	"context"

	"erp-job/internal/config"
	"erp-job/internal/database"
)

func RunMigrations(ctx context.Context, configPath string) error {
	databaseConfig, err := config.LoadDatabase(configPath)
	if err != nil {
		return err
	}

	db, err := database.OpenMySQL(databaseConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	return database.ApplyMigrations(ctx, db)
}
