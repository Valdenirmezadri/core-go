package slices

// Filter retorna um novo slice contendo apenas os elementos do slice de entrada que satisfazem o predicado fornecido.
func Filter[T any](input []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}
