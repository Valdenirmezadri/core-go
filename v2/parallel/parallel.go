package parallel

import (
	"context"
	"slices"
	"sync"
)

// Run executa as funções em paralelo e retorna todos os erros encontrados.
func Run(fns ...func() error) []error {
	raw := make([]error, len(fns))
	var wg sync.WaitGroup

	wg.Add(len(fns))
	for i, fn := range fns {
		go func() {
			defer wg.Done()
			raw[i] = fn()
		}()
	}
	wg.Wait()

	return slices.DeleteFunc(raw, func(err error) bool { return err == nil })
}

// RunWithContext executa as funções em paralelo respeitando o contexto.
// Se o contexto for cancelado, as funções que ainda não iniciaram são ignoradas.
func RunWithContext(ctx context.Context, fns ...func(ctx context.Context) error) []error {
	wrapped := make([]func() error, len(fns))
	for i, fn := range fns {
		wrapped[i] = func() error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fn(ctx)
		}
	}
	return Run(wrapped...)
}
