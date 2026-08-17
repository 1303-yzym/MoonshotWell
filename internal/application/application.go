package application

import (
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
)

type Application struct {
	Service *service.Service
}

// InitApplication 初始化应用层.
func InitApplication(appState *state.AppState) *Application {
	appService := service.New(appState)

	//// 初始化队列and消费者
	//if err := consumer.New(appState, appService); err != nil {
	//	logger.App().Fatal("InitQueue", zap.Error(err))
	//}
	//
	//// 定时任务
	//err := worker.AddJobs(appState.Cron, worker.Jobs(appService)...)
	//if err != nil {
	//	logger.App().Fatal("failed to add scheduled task", zap.Error(err))
	//}
	//
	//appState.Cron.Start()

	return &Application{
		Service: appService,
	}
}
