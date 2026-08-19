package stopper

import (
	"context"
	"sync"
)

type Stoppers interface {
	Stop()
	Context() context.Context
}
type stopper struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func New() Stoppers {
	s := &stopper{}
	s.mu = sync.Mutex{}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *stopper) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
	}
}

func (s *stopper) Context() context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ctx
}
