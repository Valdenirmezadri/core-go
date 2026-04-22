package withcontext

import "context"

// InBatch splits items into chunks of size and calls fn for each chunk,
// respecting context cancellation between batches.
func InBatch[T any](ctx context.Context, items []T, size int, fn func(batch []T) error) error {
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
