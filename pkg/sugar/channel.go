package sugar

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrClosed = errors.New("subscriber closed")
	ErrFull   = errors.New("message channel full")
)

type Channel[V any] struct {
	ch chan V

	isClosed atomic.Bool
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewChannel[V any](ctx context.Context, bufSize uint) *Channel[V] {
	ctx, cancel := context.WithCancel(ctx)

	return &Channel[V]{
		ch:     make(chan V, bufSize),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Channel[V]) Channel() chan V {
	return s.ch
}

func (s *Channel[V]) Push(v V) error {
	if s.isClosed.Load() {
		return ErrClosed
	}

	select {
	case s.ch <- v:
		return nil
	case <-s.ctx.Done():
		return ErrClosed
	default:
		return ErrFull
	}
}

func (s *Channel[V]) PushWithContext(ctx context.Context, v V) error {
	if s.isClosed.Load() {
		return ErrClosed
	}

	select {
	case s.ch <- v:
		return nil
	case <-s.ctx.Done():
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Channel[V]) PushWithTimeout(msg V, timeout time.Duration) error {
	if timeout <= 0 {
		return s.Push(msg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.PushWithContext(ctx, msg)
}

func (s *Channel[V]) Close() {
	s.isClosed.Store(true)

	if s.cancel != nil {
		s.cancel()
	}
}
