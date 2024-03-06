package safe

import "sync"

func NewItemWithData[T any](data T) Item[T] {
	item := newItem[T]()
	item.Set(data)
	return item
}

func NewItem[T any]() Item[T] {
	return newItem[T]()
}

func newItem[T any]() Item[T] {
	var data T
	return &item[T]{
		data: data,
		lock: sync.RWMutex{},
	}
}

type Item[T any] interface {
	Set(c T)
	Get() T
}

type item[T any] struct {
	data T
	lock sync.RWMutex
}

func (p *item[T]) Set(c T) {
	p.lock.Lock()
	p.data = c
	p.lock.Unlock()
}

func (p *item[T]) Get() T {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.data
}
