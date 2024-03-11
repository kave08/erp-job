package logger

import (
	"erp-job/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.SugaredLogger

func Initialize() {

	logfile := fmt.Sprintf("%s/debuglog-%s.log", config.Cfg.App.LogPath, time.Now().Format("2006-01-02T15:04"))

	file, err := os.Create(logfile)
	if err != nil {
		panic(fmt.Sprintf("failed to create log file: %s", err))
	}

	// Create a WriteSyncer for the log file
	lowPriorityOutput := zapcore.AddSync(file)

	highPriorityOutput := zapcore.Lock(os.Stdout)

	stdoutEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// Use the zapcore.Level values directly for setting the log level
	core := zapcore.NewTee(
		zapcore.NewCore(stdoutEncoder, highPriorityOutput, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, lowPriorityOutput, zapcore.ErrorLevel),
	)

	lg = zap.New(core).Sugar()
}

func Logger() *zap.SugaredLogger {
	return lg
}
