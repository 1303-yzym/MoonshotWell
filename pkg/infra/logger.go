package infra

import (
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

type LogConfig struct {
	LogDir     string `mapstructure:"log_dir" json:"log_dir" comment:"日志存储目录"`
	LogLevel   string `mapstructure:"log_level" json:"log_level" comment:"日记存储级别 debug | info | warn | error"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" comment:"单日记文件的最大大小 MB"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" comment:"备份数量"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age" comment:"备份天数"`
	Compress   bool   `mapstructure:"compress" json:"compress" comment:"是否压缩"`
	LocalTime  bool   `mapstructure:"local_time" json:"local_time" comment:"使用本地时间"`
}

type ServerInfo struct {
	ServerName  string
	ServiceName string
	Version     string
	ReVersion   string
}

func (s ServerInfo) LogField() []zapcore.Field {
	return []zapcore.Field{
		zap.String("server", s.ServerName),
		zap.String("service", s.ServiceName),
		zap.String("version", s.Version),
		zap.String("re-vision", s.ReVersion),
	}
}

func InitLogger(cfg LogConfig, isDev bool, fds ...zap.Field) {
	// 阿里主机id
	//instanceId := aliSdk.GetEscInstanceId()
	//if instanceId != nil {
	//	fds = append(fds, zap.String("instance", *instanceId))
	//}

	var ll *logger.Logs
	if isDev {
		ll = initLogs(cfg, isDev)
	} else {
		ll = initLogs(cfg, isDev, fds...)
	}

	logger.Logger.Replace(ll)

	// TODO 需要接入grpc时注入logger
	// grpc log
	//grpc_zap.ReplaceGrpcLoggerV2WithVerbosity(
	//	ll.AppLog.WithOptions(
	//		zap.IncreaseLevel(zap.WarnLevel),
	//		zap.AddCallerSkip(1),
	//	), 0,
	//)
}

func initLogs(cfg LogConfig, isDev bool, fds ...zap.Field) *logger.Logs {
	logLv, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		zap.L().Error("parse log level err", zap.Error(err))
		logLv = zap.InfoLevel
	}

	var opts []logger.LogOption

	opts = append(opts,
		logger.WithLogSliceConfig(logger.LogSliceConfig{
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  cfg.LocalTime,
		}),
		logger.WithLogDir(cfg.LogDir),
	)

	if isDev {
		opts = append(opts,
			logger.WithDevelopment(true),
		)
	}

	factory := logger.NewFactory(logLv, opts...)

	ll := factory.New(false, nil, "runtime").With(fds...)

	runtimeLog := ll.Named("runtime")
	zap.ReplaceGlobals(runtimeLog)

	return &logger.Logs{
		Original:  ll,
		AppLog:    runtimeLog,
		AccessLog: factory.New(false, runtimeLog, "access").With(fds...).Named("access"),
		SqlLog:    factory.New(false, runtimeLog, "sql").With(fds...).Named("sql"),
		EventLog:  factory.New(false, runtimeLog, "event").With(fds...).Named("event"),
		ErrorLog:  factory.New(true, runtimeLog, "error").With(fds...).Named("error"),
	}
}
