package warp

import (
	"context"
	"reflect"
)

type NilReq interface {
	HasInnerNil() struct{}
}
type WrapHandlerFunc[O context.Context, C Context[O], Q IMeta, R any] func(ctx C, req Q) (resp R, err error)

func Handle[O context.Context, C Context[O], Q IMeta, R any](rg RouterGroup[O, C], handlerFunc WrapHandlerFunc[O, C, Q, R], middlewareChain ...Middleware[O, C]) {
	var req Q

	meta, err := parseMeta(req)
	if err != nil {
		panic(err)
	}

	handlersValue := reflect.ValueOf(handlerFunc)

	route := &Route[O, C]{
		MetaInfo: &MetaInfo{
			Path:    meta.Path,
			Method:  meta.Method,
			Comment: meta.Comment,
		},
		HandlerValue: &handlersValue,
		Handlers:     []HandlerFunc[O, C]{wrapHandler(handlerFunc)},
	}

	rg.registerHandler(route, middlewareChain...)
}

func wrapHandler[O context.Context, C Context[O], Q IMeta, R any](handlerFunc WrapHandlerFunc[O, C, Q, R]) HandlerFunc[O, C] {
	return func(ctx C) {
		var req Q
		// WARNING 由于存在req为struct{}的强制判定，因此此处的判定永远为false
		// // Nil无绑定 handler中自定义绑定
		// if _, ok := any(req).(Nil); !ok {
		// 	if err := ctx.BindRequest(&req); err != nil {
		// 		return
		// 	}
		// }

		// NilReq无绑定 handler中自定义绑定
		if _, ok := any(req).(NilReq); !ok {
			if err := ctx.BindRequest(&req); err != nil {
				return
			}
		}

		resp, err := handlerFunc(ctx, req)
		if err != nil {
			ctx.Fail(err)

			return
		}

		// Nil无默认返回 handler中自定义输出
		if _, ok := any(resp).(Nil); !ok {
			ctx.Results(resp)
		}

		return
	}
}
