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
	Update(fn func(T) T)
	UpdateErr(fn func(T) (T, error)) error
	Get() T
	Read(fn func(T) error) error
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

func (p *item[T]) UpdateErr(fn func(T) (T, error)) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	newData, err := fn(p.data)
	if err != nil {
		return err
	}

	p.data = newData

	return nil
}

func (p *item[T]) Update(fn func(T) T) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.data = fn(p.data)
}

func (p *item[T]) Get() T {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.data
}

func (p *item[T]) Read(fn func(data T) error) error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return fn(p.data)
}
