package signal

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignal_Singleton(t *testing.T) {
	var called atomic.Int32
	var hookNum int32 = 20

	for i := 0; i < int(hookNum); i++ {
		On(func(ctx context.Context) {
			called.Add(1)
			Stop()
		})
	}

	go func() {
		time.Sleep(5 * time.Second)
		Stop()
	}()

	Listen()

	assert.Equal(t, called.Load(), hookNum)
}
