package base

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"vpn-web.funcworks.net/gb"
)

func initLogger() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.Encoding = "console"
	cfg.DisableStacktrace = true

	l, _ := cfg.Build()
	gb.Logger = l.Sugar()
}
