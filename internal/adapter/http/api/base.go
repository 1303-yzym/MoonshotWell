package api

import (
	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
)

// BaseHandler @Group("/base").
type BaseHandler interface {
	Health(ctx *gc.Context, req HealthReq) (HealthRes, error)
}

type (
	HealthReq struct {
		warp.Meta `path:"health" method:"post" comment:"健康检查"`
	}

	HealthRes struct {
	}
)
