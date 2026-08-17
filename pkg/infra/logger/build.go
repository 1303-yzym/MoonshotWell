package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/fatih/color"
)

type LogSliceConfig struct {
	MaxSize    int  `json:"maxsize" toml:"max_size"`
	MaxAge     int  `json:"max_age" toml:"max_age"`
	MaxBackups int  `json:"max_backups" toml:"max_backups"`
	LocalTime  bool `json:"localtime" toml:"local_time"`
	Compress   bool `json:"compress" toml:"compress"`
}

type LogFactory struct {
	Development    bool
	Level          zapcore.Level
	LogDir         string
	logSliceConfig LogSliceConfig
}

type LogOption func(*LogFactory)

func WithLogDir(logDir string) LogOption {
	return func(factory *LogFactory) {
		factory.LogDir = logDir
	}
}

func WithLogSliceConfig(logSliceConfig LogSliceConfig) LogOption {
	return func(factory *LogFactory) {
		factory.logSliceConfig = logSliceConfig
	}
}

func WithDevelopment(development bool) LogOption {
	return func(factory *LogFactory) {
		factory.Development = development
	}
}

func DefaultLogSliceConfig() LogSliceConfig {
	return LogSliceConfig{
		MaxSize:    50,
		MaxAge:     30,
		MaxBackups: 50,
		LocalTime:  true,
		Compress:   true,
	}
}

func InitBasicLog() *zap.Logger {
	factory := NewFactory(zapcore.InfoLevel, WithDevelopment(true))
	ll := factory.Basic()
	zap.ReplaceGlobals(ll)

	return ll
}

func NewFactory(lv zapcore.Level, opts ...LogOption) *LogFactory {
	l := &LogFactory{
		Development:    false,
		Level:          lv,
		LogDir:         "./logs",
		logSliceConfig: DefaultLogSliceConfig(),
	}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *LogFactory) New(logName string, fls ...zap.Field) *zap.Logger {
	var outs []*OutLogger
	if l.Development {
		outs = append(outs, &OutLogger{
			Encoder:      developmentEncoder(),
			Writer:       os.Stdout,
			LevelEnabler: l.Level,
		})
	} else {
		if !strings.HasSuffix(logName, ".log") {
			logName = logName + ".log"
		}

		outs = append(outs, &OutLogger{
			Encoder: productionEncoder(),
			Writer: &lumberjack.Logger{
				Filename:   filepath.Join(l.LogDir, logName),
				LocalTime:  l.logSliceConfig.LocalTime,
				MaxSize:    l.logSliceConfig.MaxSize,
				MaxAge:     l.logSliceConfig.MaxAge,
				MaxBackups: l.logSliceConfig.MaxBackups,
				Compress:   l.logSliceConfig.Compress,
			},
			LevelEnabler: l.Level,
		})
	}

	log := newLog(outs...)

	return log.With(fls...)
}

// Basic 基本.
func (l *LogFactory) Basic() *zap.Logger {
	var encoder zapcore.Encoder
	if l.Development {
		encoder = developmentEncoder()
	} else {
		encoder = productionEncoder()
	}

	return newLog(&OutLogger{
		Encoder:      encoder,
		Writer:       os.Stdout,
		LevelEnabler: l.Level,
	})
}

type OutLogger struct {
	Encoder      zapcore.Encoder
	Writer       io.Writer
	LevelEnabler zapcore.Level
}

func newLog(writers ...*OutLogger) *zap.Logger {
	if len(writers) < 1 {
		return nil
	}

	var logger *zap.Logger

	z := zap.NewProductionEncoderConfig()
	z.EncodeTime = zapcore.ISO8601TimeEncoder

	var core zapcore.Core
	if len(writers) == 1 {
		core = zapcore.NewCore(writers[0].Encoder,
			zapcore.AddSync(writers[0].Writer), writers[0].LevelEnabler)
	} else {
		var cores []zapcore.Core
		for i := 0; i < len(writers); i++ {
			cores = append(cores, zapcore.NewCore(writers[i].Encoder,
				zapcore.AddSync(writers[i].Writer), writers[i].LevelEnabler))
		}

		core = zapcore.NewTee(cores...)
	}

	logger = zap.New(core, zap.AddCaller())
	_ = logger.Sync()

	return logger
}

func productionEncoder() zapcore.Encoder {
	z := zap.NewProductionEncoderConfig()
	z.EncodeTime = zapcore.ISO8601TimeEncoder
	z.NameKey = "type"

	return zapcore.NewJSONEncoder(z)
}

func developmentEncoder() zapcore.Encoder {
	var _levelToColor = map[zapcore.Level]func(a ...interface{}) string{
		zap.DebugLevel:  color.New(color.FgHiMagenta).Add(color.Bold).SprintFunc(),
		zap.InfoLevel:   color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
		zap.WarnLevel:   color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
		zap.ErrorLevel:  color.New(color.FgRed).Add(color.Bold).SprintFunc(),
		zap.DPanicLevel: color.New(color.BgRed).Add(color.Bold).SprintFunc(),
		zap.PanicLevel:  color.New(color.BgRed).Add(color.Bold).SprintFunc(),
		zap.FatalLevel:  color.New(color.BgRed).Add(color.Bold).SprintFunc(),
	}

	encoderConf := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		NameKey:       "type",
		CallerKey:     "caller_line",
		LevelKey:      "level_name",
		MessageKey:    "msg",
		TimeKey:       "ts",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			s := _levelToColor[level]("[" + level.CapitalString() + "]")
			enc.AppendString(s)
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			s := color.New(color.FgHiBlack).Add(color.Bold).Sprint("[" + t.Format("2006-01-02 15:04:05") + "]")
			enc.AppendString(s)
		},
		EncodeName: func(named string, enc zapcore.PrimitiveArrayEncoder) {
			c := "<" + named + ">"

			minWidth := 10
			if len(c) > minWidth {
				minWidth = len(c)
			}

			sprintf := c + strings.Repeat(" ", minWidth-len(c))
			enc.AppendString(color.New(color.FgGreen).Add(color.Bold).Sprint(sprintf))
		},
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	})

	return encoderConf
}
