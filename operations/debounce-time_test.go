package operations

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounce_ExecutesOnlyLast(t *testing.T) {
	debounce := Debounce{}.New(100 * time.Millisecond)
	ctx := context.Background()

	var count int32

	// Chama Next várias vezes rapidamente
	for i := 0; i < 5; i++ {
		go func() {
			_ = debounce.Next(ctx, func() error {
				atomic.AddInt32(&count, 1)
				return nil
			})
		}()
		time.Sleep(20 * time.Millisecond)
	}

	// Espera tempo suficiente para o debounce disparar
	time.Sleep(200 * time.Millisecond)

	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("esperado executar apenas 1 vez, executou %d", count)
	}
}

func TestDebounce_ContextCancel(t *testing.T) {
	debounce := Debounce{}.New(200 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- debounce.Next(ctx, func() error {
			return errors.New("não deveria executar")
		})
	}()

	// Cancela o contexto antes do tempo do debounce
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Errorf("esperado erro de contexto cancelado, veio nil")
		}
	case <-time.After(300 * time.Millisecond):
		t.Errorf("timeout esperando resultado do debounce")
	}
}

func TestDebounce_ExecutesAfterWait(t *testing.T) {
	debounce := Debounce{}.New(50 * time.Millisecond)
	ctx := context.Background()

	called := make(chan struct{}, 1)
	go func() {
		_ = debounce.Next(ctx, func() error {
			called <- struct{}{}
			return nil
		})
	}()

	select {
	case <-called:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Errorf("função não executada após o tempo de espera")
	}
}
