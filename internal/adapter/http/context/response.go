package gc

import (
	"fmt"
	"net/http"

	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/config"
	"github.com/1303-yzym/MoonshotWell/pkg/code"
	"github.com/1303-yzym/MoonshotWell/pkg/i18n"
	"go.uber.org/zap/zapcore"
)

// ResponseWarp 响应结构包装.
type ResponseWarp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"error,omitempty"` // 具体错误信息 (开发环境)
}

func (r ResponseWarp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("code", r.Code)
	enc.AddString("msg", r.Msg)

	if r.Errors == nil {
		return nil
	}

	if err := enc.AddReflected("errors", r.Errors); err != nil {
		enc.AddString("errors", fmt.Sprintf("zap objectencoder reflected err: %v, value type: %T", err, r.Errors))
	}

	return nil
}

func response(ctx *Context, code code.Code, data any) {
	// i18n
	lang := ctx.G().GetHeader("Accept-Language")
	if lang == "" {
		lang = i18n.Default
	}

	// decode
	httpCode, co, msg, respErr := code.Decode(i18n.M.Load(), lang)

	var responseWarp ResponseWarp
	if co == 2000 {
		responseWarp = ResponseWarp{
			Code: co,
			Msg:  msg,
			Data: data,
		}
		ctx.G().PureJSON(httpCode, responseWarp)
	} else {
		if httpCode == 0 {
			httpCode = http.StatusBadRequest
		}

		responseWarp = ResponseWarp{
			Code: co,
			Msg:  msg,
			Data: data,
		}

		if err, ok := respErr.(error); ok {
			respErr = err.Error()
		}

		// todo 这个应该从上下文取
		if config.Load().IsDev() {
			responseWarp.Errors = respErr
		}

		ctx.G().AbortWithStatusJSON(httpCode, responseWarp)
	}

	ctx.SetResponse(responseWarp)
}

func Response(ctx *Context, code code.Code, data any) {
	response(ctx, code, data)

	return
}
