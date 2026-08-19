package semaphores

import "sync"

/*
SimpleV2 serializa trabalho por chave: uma conversa, um arquivo, um pedido.

Diferenças em relação ao Simple, que motivaram a versão nova:

  - a criação da trava é atômica. No Simple é um get seguido de um set sem nada
    segurando entre os dois, então duas goroutines na mesma chave podem receber
    mutexes diferentes e travar cada uma o seu, sem serializar nada.

  - a entrada não expira por tempo. No Simple ela expira por TTL, e quem
    segurasse a trava por mais tempo que o TTL perdia a exclusão mútua: o
    próximo Get criava outro mutex e entrava junto.

  - a limpeza é exata. Em vez de tempo, conta quem está usando e remove a
    entrada quando o último solta — então a chave pode ser sem teto, como um ID
    de conversa, sem o mapa crescer para sempre.

Devolve a função que destrava em vez do mutex, e é isso que permite contar: sem
esse retorno o pacote não teria como saber quando quem pegou terminou.
*/
type SimpleV2[K comparable] interface {
	// Lock trava a chave e devolve a função que destrava. Chamar o unlock mais
	// de uma vez é inofensivo
	Lock(key K) (unlock func())
}

func NewByKeyV2[K comparable]() SimpleV2[K] {
	return &simpleV2[K]{locks: map[K]*entry{}}
}

// entry é o mutex de uma chave junto da contagem de quem o está usando
type entry struct {
	mu   sync.Mutex
	refs int
}

type simpleV2[K comparable] struct {
	// mx protege o mapa, para todos na mesma chave pegarem o mesmo mutex
	mx    sync.Mutex
	locks map[K]*entry
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

// acquire pega ou cria a entrada e registra mais um interessado nela
func (s *simpleV2[K]) acquire(key K) *entry {
	s.mx.Lock()
	defer s.mx.Unlock()

	e, ok := s.locks[key]
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

	// confere que ainda é a mesma entrada, para não apagar uma criada depois
	if atual, ok := s.locks[key]; ok && atual == e {
		delete(s.locks, key)
	}
}
