package warp

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type ginContext struct {
	*gin.Context
}

var _ Context[*gin.Context] = (*ginContext)(nil)

func (t *ginContext) G() *gin.Context {
	return t.Context
}

func (t *ginContext) BindRequest(req any) error {
	return t.G().ShouldBind(req)
}

func (t *ginContext) Fail(err error) {
	t.G().PureJSON(http.StatusBadRequest, map[string]interface{}{
		"err": err,
	})
}

func (t *ginContext) Results(result any) {
	t.G().PureJSON(http.StatusOK, result)
}

type ginRouter struct {
}

type (
	HealthReq struct {
		Meta `path:"health" method:"post" comment:"health check"`
		Sin  int `json:"sin"`
	}
	HealthResp struct {
		Status string `json:"status"`
	}
)

func (b ginRouter) Health(ctx *ginContext, req HealthReq) (resp HealthResp, err error) {
	fmt.Println("Handle")
	fmt.Println(ctx.G().MustGet("trackId"))

	return HealthResp{Status: "ok"}, nil
}

func TestGin(t *testing.T) {
	hl := ginRouter{}

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	gwp := GinWarp(g, func(ctx *gin.Context) *ginContext {
		trackId := ulid.Make().String()
		ctx.Set("trackId", trackId)
		ctx.Header("trackId", trackId)

		return &ginContext{ctx}
	}, WithNoRouterInfo(false))

	gwp.Use(WrapMiddleware("a", func(ctx *ginContext) {
		fmt.Println("middleware 1")
	}))

	gr := gwp.Group("/base", WrapMiddleware("b", func(ctx *ginContext) {
		fmt.Println("middleware 2")
	}))

	gr2 := gr.Group("/1", WrapMiddleware("c", func(ctx *ginContext) {
		fmt.Println("middleware 3")
	}))

	Handle(gr2, hl.Health, WrapMiddleware("d", func(ctx *ginContext) {
		fmt.Println("middleware 4")
	}))

	gr2.Handle("POST", "/test", func(ctx *ginContext) {

	})

	srv := gwp.Server("0.0.0.0", "12931", nil)

	if err := func() error {
		var err error
		go func() {
			err = srv.ListenAndServe()
		}()

		return err
	}(); err != nil {
		t.Error(err.Error())
		t.Fatal("fail to run server")
	}

	time.Sleep(3 * time.Second)
}
