package service

import (
	"context"

	"github.com/1303-yzym/MoonshotWell/internal/application/repository"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
)

type ExampleService struct {
	appState   *state.AppState
	repository *repository.Repository
}

func (e ExampleService) SayHello(ctx context.Context) (err error) {
	logger.App().Info("你好！！！！")

	return
}
