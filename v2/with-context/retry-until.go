package withcontext

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type OptFunc func(*opts)

type opts struct {
	timeOutSec  uint
	MaxAttempt  uint
	delayBetwen uint
	//um sistema que vai aumentando o tempo entre tentativas
	delayMore bool
}

func (OptFunc) DelayMore(o *opts) { o.delayMore = true }

func (OptFunc) TimeOutSec(sec uint) OptFunc { return func(o *opts) { o.timeOutSec = sec } }

func (OptFunc) DelayBetwen(sec uint) OptFunc { return func(o *opts) { o.delayBetwen = sec } }

func (OptFunc) Attempts(max uint) OptFunc { return func(o *opts) { o.MaxAttempt = max } }

func defaultOpts() opts { return opts{} }

// ErrRetry signals that the operation should be retried.
var ErrRetryAgain = errors.New("retry")

// ErrLimitAttempts signals that RetryUntil gave up after Attempts(n) tries.
var ErrLimitAttempts = errors.New("attempt limit reached")

// RetryUntil loops fn until it returns something other than ErrRetry,
// or until ctx is cancelled. Use ErrRetry inside fn to signal "try again".
//
// Attempts(n) caps the number of tries; the default of zero keeps looping until
// fn stops asking or ctx dies, which is what callers retrying a collision want.
//
// Works with any return type T:
//   - (models.Propaga, error) → RetryUntil[models.Propaga](...)
//   - (bool, error)           → RetryUntil[bool](...)
//   - ([]Foo, error)          → RetryUntil[[]Foo](...)
//   - error only              → RetryUntil[struct{}](...)
func RetryUntil[T any](ctx context.Context, fn func() (T, error), opts ...OptFunc) (T, error) {
	o := defaultOpts()

	// apply, e não fn: o parâmetro já se chama fn, e sombreá-lo aqui deixaria a
	// próxima linha que usasse fn dentro do laço chamando a opção sem querer
	for _, apply := range opts {
		apply(&o)
	}

	var zero T
	var attempts uint

	for {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()

		default:
			result, err := fn()
			if err == nil {
				return result, nil
			}

			attempts++

			/*
				Teto atingido: devolve o que a última tentativa deu, incluindo o erro.

				O erro é o próprio ErrRetryAgain, e é por isso que vale envolvê-lo em fn
				com o motivo real — errors.Is atravessa o wrap, então o sinal continua
				valendo e quem recebe fica sabendo por que desistimos.
			*/
			if o.MaxAttempt > 0 && attempts >= o.MaxAttempt {
				return zero, fmt.Errorf(`%w "%d": %w`, ErrLimitAttempts, o.MaxAttempt, err)
			}

			/*
				A espera vem depois do teto: chegando ao limite não há próxima tentativa,
				e dormir ali só atrasaria um erro que já está decidido.
			*/
			if err := wait(ctx, o.delayBetwen); err != nil {
				var zero T

				return zero, err
			}

			if !errors.Is(err, ErrRetryAgain) {
				return result, err
			}
		}
	}
}

/*
wait segura a próxima tentativa por sec segundos.

Pelo select, e não por time.Sleep: com o contexto cancelado a goroutine ficaria
presa até o fim da espera, e quem cancelou já não quer o resultado.

O timer é parado no fim porque a espera acaba antes dele sempre que o contexto
morre — sem isso ele sobreviveria até disparar sozinho.

Zero não espera nada, que é o padrão de quem só quer repetir até dar certo.
*/
func wait(ctx context.Context, sec uint) error {
	if sec == 0 {
		return nil
	}

	timer := time.NewTimer(time.Duration(sec) * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		return nil
	}
}
