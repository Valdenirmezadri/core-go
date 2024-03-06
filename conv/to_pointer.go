package conv

func Pointer[T any](d T) *T {
	return &d
}
