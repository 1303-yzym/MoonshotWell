package warp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginAdapter[C Context[*gin.Context]] struct {
	engine *gin.Engine

	wrapContext func(c *gin.Context) C
}

var _ Adapter[*gin.Context, Context[*gin.Context]] = (*ginAdapter[Context[*gin.Context]])(nil)

func (g *ginAdapter[C]) WrapContext(original *gin.Context) C {
	return g.wrapContext(original)
}

func (g *ginAdapter[C]) RegisterHandler(route *Route[*gin.Context, C]) {
	var handlers []gin.HandlerFunc

	for _, handler := range route.Handlers {
		handlers = append(handlers, func(ctx *gin.Context) {
			handler(g.wrapContext(ctx))
		})
	}

	g.engine.Handle(route.Method, route.Path, handlers...)
}

func (g *ginAdapter[C]) Handler() http.Handler {
	return g.engine.Handler()
}

func GinWarp[C Context[*gin.Context]](engine *gin.Engine, wrapContext func(c *gin.Context) C, opts ...Options) *Warp[*gin.Context, C] {
	return New[*gin.Context, C](&ginAdapter[C]{
		engine:      engine,
		wrapContext: wrapContext,
	}, opts...)
}
