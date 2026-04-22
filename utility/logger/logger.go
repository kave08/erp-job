package logger

import (
	"erp-job/config"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	lg       *zap.SugaredLogger
	initOnce sync.Once
)

func Initialize() {
	initOnce.Do(func() {
		lg = buildLogger(config.Cfg.App.LogPath)
	})
}

func buildLogger(logPath string) *zap.SugaredLogger {
	highPriorityOutput := zapcore.Lock(os.Stdout)
	stdoutEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	cores := []zapcore.Core{
		zapcore.NewCore(stdoutEncoder, highPriorityOutput, zapcore.DebugLevel),
	}

	if logPath != "" {
		if err := os.MkdirAll(logPath, 0o755); err == nil {
			logfile := filepath.Join(logPath, fmt.Sprintf("debuglog-%s.log", time.Now().Format("2006-01-02T15:04")))
			if file, err := os.Create(logfile); err == nil {
				fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
				cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapcore.ErrorLevel))
			}
		}
	}

	return zap.New(zapcore.NewTee(cores...)).Sugar()
}

func Logger() *zap.SugaredLogger {
	if lg == nil {
		lg = buildLogger("")
	}

	return lg
}
