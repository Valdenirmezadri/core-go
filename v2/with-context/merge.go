package withcontext

import "context"

// Merge returns a new context that is cancelled when either parent or other is cancelled,
// but with different scopes:
//
//   - If parent is cancelled, all contexts created via Merge with the same parent are cancelled.
//   - If other is cancelled, only this specific child context is cancelled.
//
// This is useful when a long-lived owner (e.g. a handler) holds parent, and short-lived
// callers each pass their own other. Stopping the owner cancels everything; stopping a
// single caller cancels only its own work.
//
// The returned context inherits the values and deadline of parent.
// A background goroutine is started to watch other; it exits as soon as either context is done.
func Merge(parent, other context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-other.Done():
			cancel()
		case <-ctx.Done():
			// parent cancelled
		}
	}()
	return ctx
}
