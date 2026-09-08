package withcontext

import "context"

/*
DefaultBatchSize is used when the caller passes a non-positive size.

A thousand keeps well under Postgres' 65535 bind-parameter ceiling, where each id
of an IN clause costs one parameter, while still resolving the common case in a
single round trip.
*/
const DefaultBatchSize = 1000

// InBatch splits items into chunks of size and calls fn for each chunk,
// respecting context cancellation between batches.
//
// A non-positive size falls back to DefaultBatchSize: advancing by zero would
// loop forever instead of failing.
func InBatch[T any](ctx context.Context, items []T, size int, fn func(batch []T) error) error {
	if size <= 0 {
		size = DefaultBatchSize
	}

	for i := 0; i < len(items); i += size {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(items[i:min(i+size, len(items))]); err != nil {
				return err
			}
		}
	}
	return nil
}
