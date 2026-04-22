package app

import (
	"context"
	"time"

	"erp-job/internal/config"
	"erp-job/internal/database"
	"erp-job/internal/logging"
	"erp-job/internal/observability"
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
	defer func() { _ = logger.Sync() }()

	runID, err := observability.NewRunID()
	if err != nil {
		return err
	}
	ctx = observability.WithRunID(ctx, runID)
	logger = logger.With("run_id", runID)

	telemetry, shutdownTelemetry, err := observability.New(ctx, cfg.OTel)
	if err != nil {
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := shutdownTelemetry(shutdownCtx); err != nil {
			logger.Warnw("telemetry shutdown failed", "error", err.Error())
		}
	}()

	db, err := database.OpenMySQL(cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := database.ApplyMigrations(ctx, db); err != nil {
		return err
	}

	checkpoints := mysqlstore.New(db)
	source := sourcefararavand.NewClient(cfg.FararavandApp, telemetry, logger)
	target := targetaryan.NewClient(cfg.AryanApp, telemetry, logger)

	job := transfer.NewJob(checkpoints, source, target, logger, telemetry)
	return job.Run(ctx)
}
