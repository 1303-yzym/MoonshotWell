package router

import (
	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/internal/application/handler"
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"github.com/gin-gonic/gin"
)

type InitRouter[C warp.Context[*gin.Context]] interface {
	InitRouter(gr warp.RouterGroup[*gin.Context, *gc.Context])
}

func registerRoutingGroup(gr warp.RouterGroup[*gin.Context, *gc.Context], rs ...InitRouter[warp.Context[*gin.Context]]) {
	for _, r := range rs {
		r.InitRouter(gr)
	}
}

func SetupRouter(gr warp.RouterGroup[*gin.Context, *gc.Context], appState *state.AppState, appService *service.Service) {
	// 在这创建全局的处理器
	hdr := handler.New(appState, appService)

	registerRoutingGroup(
		gr,
		Base{hdr},
	)
}
