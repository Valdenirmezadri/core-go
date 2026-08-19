package slices

import "slices"

func Contains[T comparable](input []T, target T) bool {
	return slices.Contains(input, target)
}
