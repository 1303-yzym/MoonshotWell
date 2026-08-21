package contract

import (
	"context"

	"go.uber.org/zap"
)

const (
	KeyTraceId = "traceId"
)

type (
	ctxKeyTraceId struct{}
	ctxKeyLogger  struct{}
)

var (
	CtxTraceIdKey = ctxKeyTraceId{}
	CtxLoggerKey  = ctxKeyLogger{}
)

type Trace interface {
	TraceId() string
}

func TraceId(ctx context.Context) (string, bool) {
	var traceId string

	traceCtx, ok := ctx.(Trace)
	if ok {
		traceId = traceCtx.TraceId()

		return traceId, true
	}

	traceId, ok = ctx.Value(CtxTraceIdKey).(string)
	if ok {
		return traceId, true
	}

	return traceId, false
}

type Logger interface {
	LOG(ctx context.Context) *zap.Logger
}

func log(ctx context.Context) *zap.Logger {
	ll, ok := ctx.Value(CtxLoggerKey).(*zap.Logger)
	if ok {
		return ll
	}

	return nil
}

func LOGWith(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if ll := log(ctx); ll != nil {
		return ll
	}

	return logger
}

func LOG(ctx context.Context) *zap.Logger {
	return LOGWith(ctx, zap.L())
}

var _ context.Context = (*ctrContext)(nil)

type ctrContext struct {
	context.Context
}

func (c *ctrContext) LOG(_ context.Context) *zap.Logger {
	logger, ok := c.Context.Value(CtxLoggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}

	return logger
}

type ContextOptions func(*ctrContext)

// CloneContext 返回新的context但是会克隆必要字段
func CloneContext(oldCtx context.Context, opts ...ContextOptions) context.Context {
	newCtx := &ctrContext{context.Background()}

	if logger := log(oldCtx); logger != nil {
		WithContextLogger(logger)(newCtx)
	}

	if traceId, ok := TraceId(oldCtx); ok {
		WithContextTraceId(traceId)(newCtx)
	}

	for _, opt := range opts {
		opt(newCtx)
	}

	return newCtx
}

// WrapContext 包裹上下文
func WrapContext(ctx context.Context, opts ...ContextOptions) context.Context {
	var newCtx *ctrContext
	if c, ok := ctx.(*ctrContext); ok {
		newCtx = c
	} else {
		newCtx = &ctrContext{ctx}
	}

	for _, opt := range opts {
		opt(newCtx)
	}

	return newCtx
}

func WithContextValue(key, value any) ContextOptions {
	return func(ctx *ctrContext) {
		ctx.Context = context.WithValue(ctx.Context, key, value)
	}
}

func WithContextLogger(logger *zap.Logger, fields ...zap.Field) ContextOptions {
	return WithContextValue(CtxLoggerKey, logger.With(fields...))
}

func WithContextTraceId(traceId string) ContextOptions {
	return WithContextValue(CtxTraceIdKey, traceId)
}

func WithContextLoggerFields(fields ...zap.Field) ContextOptions {
	return func(ctx *ctrContext) {
		if logger, ok := ctx.Value(CtxLoggerKey).(*zap.Logger); ok {
			ctx.Context = context.WithValue(ctx.Context, CtxLoggerKey, logger.With(fields...))
		}
	}
}
