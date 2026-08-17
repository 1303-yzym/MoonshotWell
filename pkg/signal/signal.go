package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type Hook func(ctx context.Context)

type Signal struct {
	hooks []Hook
	mu    sync.RWMutex

	isClose atomic.Bool

	closeCtx    context.Context
	closeCancel context.CancelFunc

	signal chan os.Signal
}

var instance atomic.Pointer[Signal]

func init() {
	instance.Store(NewSignal())
}

func Listen() {
	instance.Load().Listen()
}

func On(hook Hook) {
	instance.Load().On(hook)
}

func Stop() {
	instance.Load().Stop()
}

func NewSignal() *Signal {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	si := &Signal{
		hooks:       make([]Hook, 0),
		mu:          sync.RWMutex{},
		isClose:     atomic.Bool{},
		closeCtx:    ctx,
		closeCancel: cancel,
		signal:      sig,
	}

	go si.listen()

	return si
}

func (s *Signal) listen() {
	<-s.signal
	s.isClose.Store(true)
	s.mu.RLock()

	for _, hook := range s.hooks {
		hook(s.closeCtx)
	}

	s.mu.RUnlock()
	s.closeCancel()
}

func (s *Signal) Listen() {
	<-s.closeCtx.Done()
}

func (s *Signal) On(hook Hook) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hooks = append(s.hooks, hook)
}

func (s *Signal) Stop() {
	if !s.isClose.Load() {
		s.signal <- syscall.SIGINT
	}
}
