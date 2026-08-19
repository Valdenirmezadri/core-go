package operations

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/safe"
)

type Debounce struct {
	wait   safe.Item[time.Duration]
	cancel safe.Item[context.CancelFunc]
}

func (Debounce) New(wait time.Duration) *Debounce {
	return &Debounce{
		wait:   safe.NewItemWithData(wait),
		cancel: safe.NewItem[context.CancelFunc](),
	}
}

func (d *Debounce) Next(ctx context.Context, fn func() error) error {
	return d.run(ctx, fn)
}

func (d *Debounce) run(ctx context.Context, fn func() error) error {
	ctx, cancel := context.WithCancel(ctx)

	d.cancel.Update(func(cf context.CancelFunc) context.CancelFunc {
		if cf != nil {
			cf()
		}
		return cancel
	})

	timer := time.NewTimer(d.wait.Get())
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return fn()
	}
}
