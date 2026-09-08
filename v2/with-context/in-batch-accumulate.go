package withcontext

import "context"

// InBatchAccumulate splits items into chunks of size, calls fn for each chunk,
// and accumulates all returned slices into a single result.
//
// On error the accumulated result is dropped: half the rows returned alongside
// an error read as the whole set to any caller that only checks the slice.
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

	if err != nil {
		return nil, err
	}

	return result, nil
}
