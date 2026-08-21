package zapgorm

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type ContextFn func(ctx context.Context) []zap.Field

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	ContextFn                 ContextFn
}

type Config struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	ContextFn                 ContextFn
}

type Options func(*Config)

func WithLogLevel(level gormlogger.LogLevel) func(*Config) {
	return func(cfg *Config) {
		cfg.LogLevel = level
	}
}

func WithSlowThreshold(slowThreshold time.Duration) func(*Config) {
	return func(cfg *Config) {
		cfg.SlowThreshold = slowThreshold
	}
}

func WithSkipCallerLookup(skipCallerLookup bool) func(*Config) {
	return func(cfg *Config) {
		cfg.SkipCallerLookup = skipCallerLookup
	}
}

func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) func(*Config) {
	return func(cfg *Config) {
		cfg.IgnoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

func WithContextFn(ctxFn ContextFn) func(*Config) {
	return func(cfg *Config) {
		cfg.ContextFn = ctxFn
	}
}

func New(zapLogger *zap.Logger, opts ...Options) Logger {
	cfg := &Config{
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             200 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		ContextFn:                 nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return Logger{
		ZapLogger:                 zapLogger,
		LogLevel:                  cfg.LogLevel,
		SlowThreshold:             cfg.SlowThreshold,
		SkipCallerLookup:          cfg.SkipCallerLookup,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ContextFn:                 cfg.ContextFn,
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		ContextFn:                 l.ContextFn,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}

	l.logger(ctx).Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}

	l.logger(ctx).Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}

	l.logger(ctx).Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)

	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		logger.Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		logger.Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l Logger) withCtxLog(ctx context.Context) *zap.Logger {
	logger := l.ZapLogger
	if l.ContextFn != nil {
		fields := l.ContextFn(ctx)
		logger = logger.With(fields...)
	}

	return logger
}

func (l Logger) logger(ctx context.Context) *zap.Logger {
	logger := l.withCtxLog(ctx)

	if l.SkipCallerLookup {
		return logger
	}

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i))
		}
	}

	return logger
}
