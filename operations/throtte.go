package operations

import (
	"context"
	"time"

	"github.com/Valdenirmezadri/core-go/safe"
)

type Throttle struct {
	interval safe.Item[time.Duration]
	lastRun  safe.Item[time.Time]
}

func (Throttle) New(interval time.Duration) *Throttle {
	return &Throttle{
		interval: safe.NewItemWithData(interval),
		lastRun:  safe.NewItemWithData(time.Time{}),
	}
}

func (t *Throttle) Next(ctx context.Context, fn func() error) error {
	now := time.Now()
	var wait time.Duration

	t.lastRun.Update(func(last time.Time) time.Time {
		if last.IsZero() || now.Sub(last) >= t.interval.Get() {
			wait = 0
			return now
		}
		wait = t.interval.Get() - now.Sub(last)
		return last
	})

	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// continue
		}
	}

	// Atualiza o lastRun para agora após aguardar
	t.lastRun.Set(time.Now())
	return fn()
}
