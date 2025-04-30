package slices

// Função genérica para encontrar um elemento em um slice
func Find[T any](input []T, predicate func(T) bool) (T, bool) {
	var zero T
	for _, v := range input {
		if predicate(v) {
			return v, true
		}
	}
	return zero, false
}
