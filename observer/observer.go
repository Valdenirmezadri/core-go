package observer

import (
	"sync/atomic"

	"github.com/Valdenirmezadri/core-go/safe"
)

type Publisher[data any] interface {
	Subscribe(Listenner[data]) (ID uint32)
	Next(data)
	UnSubscribe(id uint32)
	RemoveAll()
}

type Listenner[D any] interface {
	Listen(D)
}

type publisher[D any] struct {
	async         bool
	nextHandlerID uint32
	subscribers   safe.Lister[uint32, Listenner[D]]
}

type listen[D any] struct {
	ID string
	f  func(D)
}

func NewPublisher[D any](async bool) Publisher[D] {
	return &publisher[D]{
		async:       async,
		subscribers: safe.NewList[uint32, Listenner[D]](),
	}
}

func NewListener[D any](f func(D)) Listenner[D] {
	return &listen[D]{
		f: f,
	}
}

func (l *listen[D]) Listen(data D) {
	l.f(data)
}

func (p *publisher[D]) Subscribe(o Listenner[D]) (ID uint32) {
	if o != nil {
		nextID := atomic.AddUint32(&p.nextHandlerID, 1)
		p.subscribers.Add(nextID, o)
		return nextID
	}

	return 0
}

func (p *publisher[D]) Next(data D) {
	p.subscribers.Range(func(key uint32, l Listenner[D]) bool {
		if p.async {
			go l.Listen(data)
		} else {
			l.Listen(data)
		}
		return true
	})
}

func (p *publisher[D]) UnSubscribe(id uint32) {
	p.subscribers.Remove(id)
}

func (p publisher[D]) RemoveAll() {
	ids := []uint32{}
	p.subscribers.Range(func(key uint32, value Listenner[D]) bool {
		ids = append(ids, key)
		return true
	})

	for _, id := range ids {
		p.subscribers.Remove(id)
	}
}
