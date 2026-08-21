package middleware

import (
	"net/http"
	"net/url"
	"time"

	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/pkg/algos/encode"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(isDev bool) warp.Middleware[*gin.Context, *gc.Context] {
	return warp.WrapMiddleware("logger", func(ctx *gc.Context) {
		stateTime := time.Now()

		ctx.Ctx().Next()

		cost := time.Since(stateTime)
		// 日志输出
		if isDev {
			devLog(ctx, cost.String())
		} else {
			prodLog(ctx, cost.String())
		}
	})
}

func prodLog(ctx *gc.Context, cost string) {
	gCtx := ctx.Ctx()
	query, _ := url.QueryUnescape(gCtx.Request.URL.Query().Encode())
	req, _ := ctx.GetRequest()
	resp, _ := ctx.GetResponse()

	log := ctx.Log().With(
		zap.Int("status", gCtx.Writer.Status()),
		zap.String("host", gCtx.ClientIP()),
		zap.String("method", gCtx.Request.Method),
		zap.String("path", gCtx.Request.URL.Path),
		zap.String("query", query),
		zap.String("exec_ts", cost),
		// 请求具体信息
		zap.Any("request", req),
		// 返回信息
		zap.Any("response", resp),
	)
	if gCtx.Writer.Status() == http.StatusOK || gCtx.Writer.Status() == http.StatusMovedPermanently {
		log.Info("")
	} else {
		log.Warn("")
	}
}

func devLog(ctx *gc.Context, cost string) {
	gCtx := ctx.Ctx()
	query, _ := url.QueryUnescape(gCtx.Request.URL.Query().Encode())
	msg := encode.GinMsg{
		Status: gCtx.Writer.Status(),
		Proto:  cost,
		Host:   gCtx.ClientIP(),
		Method: gCtx.Request.Method,
		Path:   gCtx.Request.URL.Path,
		Query:  query,
	}

	log := ctx.Log().WithOptions(zap.WithCaller(false))
	if gCtx.Writer.Status() == http.StatusOK || gCtx.Writer.Status() == http.StatusMovedPermanently {
		log.Sugar().Infof("%s", encode.FormatGinLog(msg))
	} else {
		log.Sugar().Warnf("%s", encode.FormatGinLog(msg))
	}
}
