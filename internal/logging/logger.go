package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logPath string) *zap.SugaredLogger {
	stdout := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapcore.DebugLevel,
	)

	cores := []zapcore.Core{stdout}

	if logPath != "" {
		if err := os.MkdirAll(logPath, 0o755); err == nil {
			logfile := filepath.Join(logPath, fmt.Sprintf("debuglog-%s.log", time.Now().Format("2006-01-02T15:04")))
			if file, err := os.Create(logfile); err == nil {
				cores = append(cores, zapcore.NewCore(
					zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
					zapcore.AddSync(file),
					zapcore.ErrorLevel,
				))
			}
		}
	}

	return zap.New(zapcore.NewTee(cores...)).Sugar()
}
