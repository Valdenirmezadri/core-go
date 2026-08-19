package safehandler

import (
	"github.com/Valdenirmezadri/core-go/safe"
)

type Item[T any] interface {
	Set(c T)
	Unset()
	Get() (T, error)
	IsSet() bool
}

type item[T any] struct {
	err      error
	data     safe.Item[T]
	hasValue safe.Item[bool]
}

func NewItem[T any](errWhenNil error) Item[T] {
	return &item[T]{
		err:      errWhenNil,
		data:     safe.NewItem[T](),
		hasValue: safe.NewItemWithData(false),
	}
}

func (p *item[T]) Unset() {
	var zero T
	p.data.Set(zero)
	p.hasValue.Set(false)
}

func (p *item[T]) IsSet() bool {
	return p.hasValue.Get()
}

func (p *item[T]) Set(value T) {
	p.data.Set(value)
	p.hasValue.Set(true)
}

func (p *item[T]) Get() (T, error) {
	if p.IsSet() {
		return p.data.Get(), nil
	}

	var zero T
	return zero, p.err
}
