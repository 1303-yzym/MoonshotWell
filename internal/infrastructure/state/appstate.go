package state

import (
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/config"
	"github.com/1303-yzym/MoonshotWell/pkg/contract"
	db "github.com/1303-yzym/MoonshotWell/pkg/infra/DB"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
)

type Version struct {
	Version  string
	Revision string
}

// AppState 全局句柄
type AppState struct {
	contract.Storage
	//
	//Version Version
	//// rabbitmq
	//RMQ *rabbitmq.RMQ
	//// 事件总线
	//EventBus *eventBus.Bus
	// 日志
	Logs *logger.Logs
	//// 配置文件
	Cfg *config.Config
	//// jwt
	//JWT jwts.JWT[jwtClaims.DcUser]
	//// 存储
	//Oss oss.Oss
	//// 第三方sdk
	//SDK *sdk.SDK
	//// 支付聚合
	//PayAggr *payAggr.PayAggr
	//// 雪花id生成器
	//Snowflake *snowflake.Snowflake
	//// 安全
	//Security *security.Security
	//// 定时任务
	//Cron *cron.Cron
}

func InitAppState() *AppState {
	// 初始化数据库
	log := logger.Logger.Load()
	cfg := config.Load()
	DB := db.InitDB(log.SqlLog, cfg.DB)
	return &AppState{
		Storage: contract.NewStorage(DB, nil, nil),
		Logs:    logger.Logger.Load(),
		Cfg:     config.Load(),
	}
}
