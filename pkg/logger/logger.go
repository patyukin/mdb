package logger

import (
	"fmt"
	"github.com/patyukin/mdb/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.Set(cfg.Logger.Level); err != nil {
		return nil, fmt.Errorf("failed to set log level: %w", err)
	}

	stdout := zapcore.AddSync(os.Stdout)
	var consoleEncoder zapcore.Encoder

	if cfg.Logger.Mode == "devel" {
		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder = zapcore.NewConsoleEncoder(developmentCfg)
	} else {
		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder = zapcore.NewConsoleEncoder(productionCfg)
	}

	return zap.New(
		zapcore.NewCore(consoleEncoder, stdout, level),
	), nil
}
