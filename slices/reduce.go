package slices

// Reduce aplica uma função de acumulação f a cada elemento do slice s, começando com o valor inicial init,
// e retornando o valor acumulado final.
func Reduce[S ~[]E, E any, R any](s S, init R, f func(R, E, int) R) R {
	acc := init
	for i, v := range s {
		acc = f(acc, v, i)
	}
	return acc
}
