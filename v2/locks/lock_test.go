package locks

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestLockSoUmPega é a prova da exclusão: cem goroutines disputando, uma entra.
func TestLockSoUmPega(t *testing.T) {
	locker := New()

	var wg sync.WaitGroup
	var pegaram int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if alreadyLocked, _ := locker.Lock(); !alreadyLocked {
				atomic.AddInt64(&pegaram, 1)
			}
		}()
	}

	wg.Wait()

	if pegaram != 1 {
		t.Fatalf("esperado 1 dono da trava, veio %d", pegaram)
	}
}

// TestReleaserDevolveATrava garante que soltar libera de verdade: sem isso a
// trava vira definitiva e quem depende dela nunca mais roda.
func TestReleaserDevolveATrava(t *testing.T) {
	locker := New()

	alreadyLocked, releaser := locker.Lock()
	if alreadyLocked {
		t.Fatal("trava nova não pode estar presa")
	}

	if outro, _ := locker.Lock(); !outro {
		t.Fatal("segundo Lock com a trava presa tinha que ser recusado")
	}

	releaser()

	if depois, _ := locker.Lock(); depois {
		t.Fatal("depois do releaser a trava tinha que estar livre")
	}
}
