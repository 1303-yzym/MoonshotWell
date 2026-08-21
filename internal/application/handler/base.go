package handler

import (
	"github.com/1303-yzym/MoonshotWell/internal/adapter/http/api"
	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
)

type BaseHandlerImpl struct {
	*service.Service
}

func (h BaseHandlerImpl) Health(ctx *gc.Context, req api.HealthReq) (resp api.HealthRes, err error) {
	return
}
