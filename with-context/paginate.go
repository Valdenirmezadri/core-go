package withcontext

import "context"

// Paginate executa paginação automática com verificação de contexto.
// fetchFn é chamado repetidamente até retornar uma lista vazia ou um erro.
// processFn é executado para cada item retornado por fetchFn.
// A paginação para quando:
//   - fetchFn retorna uma lista vazia (fim dos dados)
//   - fetchFn ou processFn retornam um erro
//   - o contexto é cancelado
func Paginate[T any](ctx context.Context, fetchFn func() ([]T, error), processFn func(item T) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			items, err := fetchFn()
			if err != nil {
				return err
			}

			if len(items) == 0 {
				return nil
			}

			if err := For(ctx, items, processFn); err != nil {
				return err
			}
		}
	}
}
