package logger

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestDev(t *testing.T) {
	factory := NewFactory(zapcore.InfoLevel,
		WithDevelopment(false),
		WithLogSliceConfig(LogSliceConfig{
			MaxSize:    50,
			MaxAge:     30,
			MaxBackups: 50,
			LocalTime:  true,
			Compress:   true,
		}),
	)
	log := factory.Basic()
	zap.ReplaceGlobals(log)
	zap.L().Named("access").Info("hello world")
	zap.L().Named("access").Info("hello world1")
	log2 := factory.New(false, nil, "rn",
		zap.String("instance", "dv-1"),
		zap.String("server", "dc-server"),
		zap.String("version", "0.1.2"),
		zap.String("env", "prod"),
	)
	zap.ReplaceGlobals(log2)
	zap.L().Info("hello world")
	access := log2.Named("access").With(zap.String("module", "service"), zap.String("traceId", "e094b8338e7f4b1894b8338e7f4b1867"))
	access.Info("access")

	// 200000
	for i := range 20 {
		access.Info("access.info", zap.Int("access.info", i))
	}

	time.Sleep(2 * time.Second)
}
