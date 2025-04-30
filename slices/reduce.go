package slices

// Reduce aplica uma função de acumulação f a cada elemento do slice s, começando com o valor inicial init,
// e retornando o valor acumulado final.
func Reduce[list ~[]E, E any, acc any](s list, init acc, f func(acc, E, int) acc) acc {
	_acc := init
	for i, v := range s {
		_acc = f(_acc, v, i)
	}
	return _acc
}
