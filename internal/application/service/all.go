package service

import (
	"github.com/1303-yzym/MoonshotWell/internal/application/repository"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
)

type Service struct {
	AppState   *state.AppState
	Repository *repository.Repository

	ExampleService ExampleService
}

func New(appState *state.AppState) *Service {
	// 初始化仓储对象
	rp := repository.New(appState.Storage, appState.Cfg.Env)
	return &Service{
		AppState:   appState,
		Repository: rp,

		ExampleService: ExampleService{appState, rp},
	}
}
