package handler

import (
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
)

type Handler struct {
	AppState *state.AppState
	// TODO 注册Handler
	// AdminHandler            api.AdminHandler
}

func New(appState *state.AppState, appService *service.Service) Handler {
	// 在这创建服务
	// sv := appService

	return Handler{
		AppState: appState,
		// AdminHandler:            NewAdminHandler(sv),

	}
}
