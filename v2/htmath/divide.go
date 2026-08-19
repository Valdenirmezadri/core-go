package htmath

import "errors"

func Divide[T Number](this, by T) (T, error) {
	var zero T
	if by == zero {
		return zero, errors.New("cannot divide by zero")
	}

	result := this / by

	return result, nil
}
