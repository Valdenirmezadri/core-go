package semaphores

import (
	"sync"
	"testing"
	"time"
)

// entradas é quantas chaves ainda estão vivas no mapa
func entradas[K comparable](t *testing.T, s SimpleV2[K]) int {
	t.Helper()

	impl, ok := s.(*simpleV2[K])
	if !ok {
		t.Fatal("SimpleV2 não é *simpleV2")
	}

	impl.mx.Lock()
	defer impl.mx.Unlock()

	return len(impl.locks)
}

// só uma goroutine por vez dentro do trecho crítico da mesma chave
func TestSimpleV2_Serializa(t *testing.T) {
	s := NewByKeyV2[string]()

	const goroutines = 50

	dentro := 0
	maximo := 0

	var guarda sync.Mutex
	var wg sync.WaitGroup

	start := make(chan struct{})

	for range goroutines {
		wg.Add(1)

		go func() {
			defer wg.Done()

			<-start

			unlock := s.Lock("chave")
			defer unlock()

			guarda.Lock()
			dentro++
			if dentro > maximo {
				maximo = dentro
			}
			guarda.Unlock()

			time.Sleep(time.Millisecond)

			guarda.Lock()
			dentro--
			guarda.Unlock()
		}()
	}

	close(start)
	wg.Wait()

	if maximo != 1 {
		t.Errorf("chegou a %d goroutines dentro do trecho crítico, queria 1", maximo)
	}
}

/*
Duas goroutines na mesma chave têm que disputar o mesmo mutex.

É o bug do Simple: criação não atômica devolvia mutexes diferentes e cada uma
travava o seu, sem serializar nada.
*/
func TestSimpleV2_MesmaChaveMesmaTrava(t *testing.T) {
	const rodadas = 2000

	for rodada := range rodadas {
		s := NewByKeyV2[string]()

		var wg sync.WaitGroup
		start := make(chan struct{})

		dentro := 0
		colidiu := false

		var guarda sync.Mutex

		for range 8 {
			wg.Add(1)

			go func() {
				defer wg.Done()

				<-start

				unlock := s.Lock("chave")
				defer unlock()

				guarda.Lock()
				dentro++
				if dentro > 1 {
					colidiu = true
				}
				guarda.Unlock()

				guarda.Lock()
				dentro--
				guarda.Unlock()
			}()
		}

		close(start)
		wg.Wait()

		if colidiu {
			t.Fatalf("rodada %d: duas goroutines entraram juntas na mesma chave", rodada)
		}
	}
}

// a entrada some quando o último solta, senão o mapa cresceria para sempre
func TestSimpleV2_LimpaQuandoNinguemSegura(t *testing.T) {
	s := NewByKeyV2[uint]()

	const chaves = 1000

	for n := range uint(chaves) {
		unlock := s.Lock(n)
		unlock()
	}

	if restantes := entradas(t, s); restantes != 0 {
		t.Errorf("sobraram %d entradas depois de todos soltarem, queria 0", restantes)
	}
}

// enquanto alguém segura ou espera, a entrada não pode sumir
func TestSimpleV2_NaoLimpaComGenteEsperando(t *testing.T) {
	s := NewByKeyV2[string]()

	unlock := s.Lock("chave")

	if restantes := entradas(t, s); restantes != 1 {
		t.Fatalf("entradas = %d com a trava na mão, queria 1", restantes)
	}

	esperando := make(chan struct{})
	soltou := make(chan struct{})

	go func() {
		close(esperando)

		segundo := s.Lock("chave")
		segundo()

		close(soltou)
	}()

	<-esperando
	time.Sleep(20 * time.Millisecond)

	// o segundo está bloqueado no mutex: a entrada segue viva pela contagem dele
	if restantes := entradas(t, s); restantes != 1 {
		t.Errorf("entradas = %d com alguém esperando, queria 1", restantes)
	}

	unlock()
	<-soltou

	if restantes := entradas(t, s); restantes != 0 {
		t.Errorf("entradas = %d depois de todos soltarem, queria 0", restantes)
	}
}

// chamar o unlock duas vezes não pode destravar o mutex de outra goroutine
func TestSimpleV2_UnlockDuasVezes(t *testing.T) {
	s := NewByKeyV2[string]()

	unlock := s.Lock("chave")
	unlock()
	unlock()

	if restantes := entradas(t, s); restantes != 0 {
		t.Errorf("entradas = %d, queria 0", restantes)
	}

	// com estado sujo, este Lock travaria para sempre ou veria contagem errada
	outro := s.Lock("chave")
	defer outro()

	if restantes := entradas(t, s); restantes != 1 {
		t.Errorf("entradas = %d depois de travar de novo, queria 1", restantes)
	}
}

// chaves diferentes não se bloqueiam
// ttl remove a entrada ociosa, mas Lock de novo antes do ttl cancela a
// remoção e reseta o relógio
func TestSimpleV2_TTLExpiraOciosaEResetaEmUso(t *testing.T) {
	const ttl = 30 * time.Millisecond
	s := NewByKeyV2[string](ttl)

	unlock := s.Lock("k")
	unlock()

	// reacessa antes do ttl: tem que cancelar a remoção agendada
	time.Sleep(ttl / 2)
	unlock = s.Lock("k")
	unlock()

	if restantes := entradas(t, s); restantes != 1 {
		t.Fatalf("entrada sumiu antes do ttl contado a partir do último unlock, restantes=%d", restantes)
	}

	time.Sleep(2 * ttl)

	if restantes := entradas(t, s); restantes != 0 {
		t.Errorf("entrada ociosa não expirou depois do ttl, restantes=%d", restantes)
	}
}

func TestSimpleV2_ChavesDiferentesNaoBloqueiam(t *testing.T) {
	s := NewByKeyV2[string]()

	primeira := s.Lock("a")
	defer primeira()

	pronto := make(chan struct{})

	go func() {
		segunda := s.Lock("b")
		segunda()

		close(pronto)
	}()

	select {
	case <-pronto:
	case <-time.After(time.Second):
		t.Error("Lock em chave diferente ficou bloqueado")
	}
}
