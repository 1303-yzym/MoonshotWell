package zapgorm

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	ll := New(zap.NewExample(), WithContextFn(func(ctx context.Context) []zap.Field {
		return []zap.Field{
			zap.String("type", "sql"),
			zap.String("server", "test"),
		}
	}))
	ll.SetAsDefault()
}
