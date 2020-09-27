package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(level string) *zap.Logger {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	switch level {
	case "DEBUG":
		atom.SetLevel(zap.DebugLevel)
	case "WARN":
		atom.SetLevel(zap.WarnLevel)
	case "ERROR":
		atom.SetLevel(zap.ErrorLevel)
	case "PANIC":
		atom.SetLevel(zap.PanicLevel)
	default:
		atom.SetLevel(zap.InfoLevel)
	}

	return logger
}