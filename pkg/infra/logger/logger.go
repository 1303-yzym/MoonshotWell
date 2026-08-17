package logger

import (
	"github.com/1303-yzym/MoonshotWell/pkg/sugar"
	"go.uber.org/zap"

	"github.com/fatih/color"
)

var (
	Logger = sugar.NewSingletonWithValue[*Logs, *zap.Logger](func(ll *zap.Logger) *Logs {
		if ll == nil {
			ll = InitBasicLog()
		}

		return &Logs{
			Original:  ll,
			AppLog:    ll,
			AccessLog: ll,
			SqlLog:    ll,
			EventLog:  ll,
		}
	})
)

func init() {
	color.NoColor = false

	Logger.Init(InitBasicLog())
}

type Logs struct {
	Original  *zap.Logger
	AppLog    *zap.Logger
	AccessLog *zap.Logger
	SqlLog    *zap.Logger
	EventLog  *zap.Logger
}

func Log() *zap.Logger {
	return Logger.Load().Original
}

func App() *zap.Logger {
	return Logger.Load().AppLog
}

func Access() *zap.Logger {
	return Logger.Load().AccessLog
}

func SQL() *zap.Logger {
	return Logger.Load().SqlLog
}

func Event() *zap.Logger {
	return Logger.Load().EventLog
}
