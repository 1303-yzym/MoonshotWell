package gc

import (
	"context"
	"errors"

	"github.com/1303-yzym/MoonshotWell/pkg/code"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
	"github.com/1303-yzym/MoonshotWell/pkg/validation"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type Context struct {
	*gin.Context
}

func New(ctx *gin.Context) *Context {
	trackId := ulid.Make().String()
	ctx.Set("trackId", trackId)
	ctx.Header("trackId", trackId)

	return &Context{Context: ctx}
}

var _ warp.Context[*gin.Context] = (*Context)(nil)

func (c *Context) G() *gin.Context {
	return c.Context
}

func (c *Context) Ctx() *gin.Context {
	return c.Context
}

func (c *Context) Fail(err error) {
	var codeErr code.Code
	if errors.As(err, &codeErr) {
		response(c, codeErr, nil)

		return
	}

	response(c, code.Err, err)
}

func (c *Context) Results(result any) {
	response(c, code.OK, result)
}

func (c *Context) BindRequest(req any) error {
	if err := c.Context.ShouldBind(req); err != nil {
		response(c, code.ErrField, validation.ParseValidationErrors(err))

		return err
	}

	return nil
}

func (c *Context) TraceId() string {
	traceId, ok := c.G().Get("traceId")
	if !ok {
		traceId = ulid.Make().String()
		c.G().Set("traceId", traceId)
	}

	return traceId.(string)
}

func (c *Context) Log() *zap.Logger {
	opts := []zap.Field{
		zap.String("traceId", c.TraceId()),
	}

	return logger.Access().With(opts...)
}

// LOG 实现 state.Logger.
func (c *Context) LOG(_ context.Context) *zap.Logger {
	return c.Log()
}

func (c *Context) GetRequest() (value any, exists bool) {
	return c.Ctx().Get("request")
}

func (c *Context) SetResponse(responseWarp ResponseWarp) {
	c.Ctx().Set("response", responseWarp)
}

func (c *Context) GetResponse() (responseWarp ResponseWarp, err error) {
	value, exists := c.Ctx().Get("response")
	if !exists {
		err = errors.New("response not found in context")

		return
	}

	responseWarp, ok := value.(ResponseWarp)
	if !ok {
		err = errors.New("response type assertion failed")

		return
	}

	return responseWarp, nil
}

//func (c *Context) MustUserInfo() state.UserInfo {
//	info, exists := c.Ctx().Get("user_info")
//	if !exists {
//		c.Log().Panic("user_info no found")
//	}
//
//	userInfo, ok := info.(state.UserInfo)
//	if !ok {
//		c.Log().Panic("user_info type error")
//	}
//
//	return userInfo
//}

//func (c *Context) MustUserId() uint64 {
//	userInfo := c.MustUserInfo()
//
//	return userInfo.UserId
//}

//func (c *Context) MustSessionId() uint64 {
//	userInfo := c.MustUserInfo()
//
//	return userInfo.UserId
//}

//func (c *Context) TryUserId() (userId uint64, isLogin bool) {
//	userInfo := c.MustUserInfo()
//
//	return userInfo.UserId, userInfo.IsLogin
//}
