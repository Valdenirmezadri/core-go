package withcontext

import (
	"context"
	"errors"
)

// ErrRetry signals that the operation should be retried.
var ErrRetryAgain = errors.New("retry")

// RetryUntil loops fn until it returns something other than ErrRetry,
// or until ctx is cancelled. Use ErrRetry inside fn to signal "try again".
//
// Works with any return type T:
//   - (models.Propaga, error) → RetryUntil[models.Propaga](...)
//   - (bool, error)           → RetryUntil[bool](...)
//   - ([]Foo, error)          → RetryUntil[[]Foo](...)
//   - error only              → RetryUntil[struct{}](...)
func RetryUntil[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	for {
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()

		default:
			result, err := fn()
			if !errors.Is(err, ErrRetryAgain) {
				return result, err
			}
		}
	}
}
