package withcontext

import "context"

// InBatchAccumulate splits items into chunks of size, calls fn for each chunk,
// and accumulates all returned slices into a single result.
func InBatchAccumulate[T any, I any](ctx context.Context, items []I, size int, fn func(batch []I) ([]T, error)) ([]T, error) {
	var result []T
	err := InBatch(ctx, items, size, func(batch []I) error {
		got, err := fn(batch)
		if err != nil {
			return err
		}
		result = append(result, got...)
		return nil
	})
	return result, err
}
