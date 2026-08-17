package http

import (
	"net/http"

	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/internal/adapter/http/router"
	"github.com/1303-yzym/MoonshotWell/internal/application"
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
	"github.com/1303-yzym/MoonshotWell/pkg/config"
	"github.com/1303-yzym/MoonshotWell/pkg/server"
	"github.com/1303-yzym/MoonshotWell/pkg/validation"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"github.com/gin-gonic/gin"
)

func initHttpSrv(appState *state.AppState, appService *service.Service) http.Handler {
	switch appState.Cfg.Env {
	case config.EnvDev:
		gin.SetMode(gin.DebugMode)
	case config.EnvTest:
		gin.SetMode(gin.TestMode)
	case config.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	}

	validation.RegisterValidations()

	gSrv := gin.New()
	gSrv.MaxMultipartMemory = 8 << 20 // 8 MiB
	//  gSrv.Use(middleware.Cors())

	// gSrv.Use(otelgin.Middleware(appState.Cfg.ServiceName))

	gSrv.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	gSrv.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "这什么都木有",
		})
	})

	gwp := warp.GinWarp(gSrv, func(ctx *gin.Context) *gc.Context {
		return gc.New(ctx)
	}, warp.WithNoRouterInfo(appState.Cfg.IsProd()))

	// gwp.Use(middleware.Logger(appState.Cfg.IsDev()))

	rootGroup := gwp.Group(appState.Cfg.Server.HTTP.RouterPrefix)

	router.SetupRouter(rootGroup, appState, appService)

	return gwp.Handler()
}

func RunHttpService(appState *state.AppState, app *application.Application) *http.Server {
	engine := initHttpSrv(appState, app.Service)

	return server.RunServer(
		"0.0.0.0",
		appState.Cfg.Server.HTTP.Port,
		engine,
		nil,
	)
}
