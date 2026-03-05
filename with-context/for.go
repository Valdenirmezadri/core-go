package withcontext

import "context"

// For itera sobre uma lista verificando cancelamento de contexto a cada item.
// Se fn retornar erro, a iteração é interrompida e o erro é retornado.
func For[T any](ctx context.Context, list []T, fn func(item T) error) error {
	for _, item := range list {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(item); err != nil {
				return err
			}
		}
	}
	return nil
}
