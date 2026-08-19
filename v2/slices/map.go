package slices

// Função genérica para mapear elementos de um slice
func Map[T any, U any](input []T, transform func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = transform(v)
	}
	return result
}

// MapToSlice transforma um map em um slice
func MapToSlice[K comparable, V any, U any](input map[K]V, transform func(K, V) U) []U {
	result := make([]U, 0, len(input))
	for k, v := range input {
		result = append(result, transform(k, v))
	}
	return result
}
