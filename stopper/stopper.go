package stopper

import (
	"context"
)

type Stopper interface {
	Stop()
	Done() <-chan struct{}
	//Context() context.Context
}
type stopper struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func New() Stopper {
	s := &stopper{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

func (s stopper) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s stopper) Done() <-chan struct{} {
	return s.ctx.Done()
}

/* func (s stopper) Context() context.Context {
	return s.ctx
}
*/
