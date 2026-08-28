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
			Original:  ll.Named("original"), // 原始日志，未经过处理的原始数据
			AppLog:    ll.Named("runtime"),  // 运行时日志，业务逻辑日志，业务流程、数据变化
			AccessLog: ll.Named("access"),   // 请求/响应日志，HTTP/RPC调用监控
			SqlLog:    ll.Named("sql"),      // 数据库日志，SQL执行、慢查询
			EventLog:  ll.Named("event"),    // 事件驱动日志，包括消息队列、事件的发布/订阅
			ErrorLog:  ll.Named("error"),    // 错误日志，需要打开错误堆栈的记录
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
	ErrorLog  *zap.Logger
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

func Error() *zap.Logger {
	return Logger.Load().ErrorLog
}
