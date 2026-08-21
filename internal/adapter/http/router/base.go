package router

import (
	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/internal/application/handler"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"github.com/gin-gonic/gin"
)

type Base struct {
	handler.Handler
}

func (c Base) InitRouter(gr warp.RouterGroup[*gin.Context, *gc.Context]) {
	hl := c.BaseHandler
	router := gr.Group("")
	{
		warp.Handle(router, hl.Health)
	}
}
