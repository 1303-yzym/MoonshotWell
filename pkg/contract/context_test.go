package contract

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	traceId := uuid.New().String()
	logger := zap.NewExample()

	tests := []struct {
		name    string
		ctxFunc func() context.Context
	}{
		{
			name: "context.ctrContext",
			ctxFunc: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "gin.ctrContext",
			ctxFunc: func() context.Context {
				req := &http.Request{}
				req.WithContext(context.Background())
				return &gin.Context{
					Request: req,
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := WrapContext(test.ctxFunc(),
				WithContextLogger(logger),
				WithContextTraceId(traceId),
			)

			tt := uuid.New().String()
			ctx = context.WithValue(ctx, "tt", tt)

			ctx = WrapContext(ctx,
				WithContextTraceId(traceId),
			)

			assert.Equal(t, tt, ctx.Value("tt"))

			assert.Equal(t, logger, LOG(ctx))
			ctxTraceId, ok := TraceId(ctx)
			assert.True(t, ok)
			assert.Equal(t, traceId, ctxTraceId)
			//
			c1 := CloneContext(ctx)
			assert.Equal(t, logger, LOG(c1))
			ctxTraceId, ok = TraceId(c1)
			assert.True(t, ok)
			assert.Equal(t, traceId, ctxTraceId)
		})
	}
}
