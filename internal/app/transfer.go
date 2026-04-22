package app

import (
	"context"

	"erp-job/internal/config"
	"erp-job/internal/database"
	"erp-job/internal/logging"
	sourcefararavand "erp-job/internal/source/fararavand"
	mysqlstore "erp-job/internal/store/mysql"
	targetaryan "erp-job/internal/target/aryan"
	"erp-job/internal/transfer"
)

func RunTransfer(ctx context.Context, configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	logger := logging.New(cfg.App.LogPath)
	db, err := database.OpenMySQL(cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	checkpoints := mysqlstore.New(db)
	source := sourcefararavand.NewClient(cfg.FararavandApp)
	target := targetaryan.NewClient(cfg.AryanApp)

	job := transfer.NewJob(checkpoints, source, target, logger)
	return job.Run(ctx)
}
