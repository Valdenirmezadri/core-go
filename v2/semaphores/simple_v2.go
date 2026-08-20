package semaphores

import (
	"sync"
	"time"
)

/*
SimpleV2 serializa trabalho por chave: uma conversa, um arquivo, um pedido.

Diferenças em relação ao Simple, que motivaram a versão nova:

  - a criação da trava é atômica. No Simple é um get seguido de um set sem nada
    segurando entre os dois, então duas goroutines na mesma chave podem receber
    mutexes diferentes e travar cada uma o seu, sem serializar nada.

  - a expiração é de ociosidade, não de tempo de vida. No Simple o TTL corre
    desde a criação, então quem segurasse a trava por mais tempo que o TTL
    perdia a exclusão mútua: o próximo Get criava outro mutex e entrava junto.
    Aqui o relógio só começa a contar quando o último solta, e qualquer Lock
    novo cancela e reseta esse relógio — então nunca expira em uso.

  - a limpeza é exata. Conta quem está usando; a entrada só é candidata a
    expirar quando o último solta, e é removida de fato se ninguém pegar de
    novo dentro do ttl.

Devolve a função que destrava em vez do mutex, e é isso que permite contar: sem
esse retorno o pacote não teria como saber quando quem pegou terminou.
*/
type SimpleV2[K comparable] interface {
	// Lock trava a chave e devolve a função que destrava. Chamar o unlock mais
	// de uma vez é inofensivo
	Lock(key K) (unlock func())
}

// NewByKeyV2 cria o semáforo. ttl, se > 0, remove uma entrada ociosa (sem
// ninguém segurando) depois desse tempo sem uso; cada Lock novo cancela a
// remoção pendente. Omitido ou 0, sem expiração — entrada só some quando o
// mapa der conta sozinho de que ninguém mais referencia a chave.
func NewByKeyV2[K comparable](ttl ...time.Duration) SimpleV2[K] {
	s := &simpleV2[K]{locks: map[K]*entry{}}
	if len(ttl) > 0 {
		s.ttl = ttl[0]
	}
	return s
}

// entry é o mutex de uma chave junto da contagem de quem o está usando
type entry struct {
	mu    sync.Mutex
	refs  int
	timer *time.Timer // remoção agendada, só existe enquanto refs == 0
}

type simpleV2[K comparable] struct {
	// mx protege o mapa, para todos na mesma chave pegarem o mesmo mutex
	mx    sync.Mutex
	locks map[K]*entry
	ttl   time.Duration
}

func (s *simpleV2[K]) Lock(key K) (unlock func()) {
	e := s.acquire(key)

	e.mu.Lock()

	var once sync.Once

	return func() {
		once.Do(func() {
			e.mu.Unlock()
			s.release(key, e)
		})
	}
}

// acquire pega ou cria a entrada, cancela remoção pendente e registra mais um
// interessado nela
func (s *simpleV2[K]) acquire(key K) *entry {
	s.mx.Lock()
	defer s.mx.Unlock()

	e, ok := s.locks[key]
	if ok && e.timer != nil {
		e.timer.Stop()
		e.timer = nil
	}
	if !ok {
		e = &entry{}
		s.locks[key] = e
	}

	/*Contado antes de disputar o mutex, e não depois de conseguir.

	Quem está esperando também precisa segurar a entrada viva: se só quem tem a
	trava contasse, a entrada seria removida no unlock e quem estava na fila
	acordaria travando um mutex já órfão, enquanto o próximo a chegar criaria
	outro e entraria junto.
	*/
	e.refs++

	return e
}

func (s *simpleV2[K]) release(key K, e *entry) {
	s.mx.Lock()
	defer s.mx.Unlock()

	e.refs--
	if e.refs > 0 {
		return
	}

	if s.ttl <= 0 {
		s.remove(key, e)
		return
	}

	// agenda a remoção; um Lock novo na chave cancela isso em acquire
	e.timer = time.AfterFunc(s.ttl, func() { s.expire(key, e) })
}

func (s *simpleV2[K]) expire(key K, e *entry) {
	s.mx.Lock()
	defer s.mx.Unlock()

	// alguém pegou a chave entre o timer disparar e esta goroutine pegar o mx
	if e.refs > 0 {
		return
	}

	s.remove(key, e)
}

// remove tira a entrada do mapa, conferindo que ainda é a mesma (não uma
// criada depois). Chamado com mx já travado.
func (s *simpleV2[K]) remove(key K, e *entry) {
	if atual, ok := s.locks[key]; ok && atual == e {
		delete(s.locks, key)
	}
}
