package adapter

import (
	"context"
	"errors"
	"net/http"

	athttp "github.com/1303-yzym/MoonshotWell/internal/adapter/http"
	"github.com/1303-yzym/MoonshotWell/internal/application"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
)

// Adapter 出口适配器
type Adapter struct {
	httpServer *http.Server
	// opsHttpServer *http.Server
	// rpcServer *grpc.Server
}

func (a *Adapter) Shutdown(ctx context.Context) error {
	var errArr error
	//if err := a.opsHttpServer.Shutdown(ctx); err != nil {
	//	errArr = errors.Join(errArr, err)
	//}

	if err := a.httpServer.Shutdown(ctx); err != nil {
		errArr = errors.Join(errArr, err)
	}

	// a.rpcServer.GracefulStop()

	return errArr
}

func LoadAdapter(appState *state.AppState, app *application.Application) *Adapter {
	// http
	srv := athttp.RunHttpService(appState, app)

	// opsSrv := athttp.RunOpsService(appState)

	return &Adapter{
		httpServer: srv,
	}
}
