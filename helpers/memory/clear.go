package memory

import (
	"context"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/Valdenirmezadri/core-go/operations"
)

type ClearMemory struct {
	debounce *operations.Throttle
}

func (ClearMemory) New(after time.Duration) *ClearMemory {
	debounce := operations.Throttle{}.New(after)
	//debounce := operations.Debounce{}.New(after)
	return &ClearMemory{
		debounce: debounce,
	}
}

func (c *ClearMemory) Run() {
	go c.debounce.Next(context.Background(), func() error {
		runtime.GC()
		debug.FreeOSMemory()
		return nil
	})
}
