package safe

import "sync"

type Lister[K comparable, V any] interface {
	Len() int
	Add(key K, value V)
	Has(key K) bool
	Load(key K) (ok bool, value V)
	Remove(key K)
	Range(f func(key K, value V) bool)
}

type list[K comparable, V any] struct {
	list sync.Map
}

func NewList[K comparable, V any]() Lister[K, V] {
	return &list[K, V]{}
}

// Len count how many stored elements in the list
func (e *list[K, V]) Len() int {
	length := 0
	e.list.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

// Add sets the value for a key.
func (l *list[K, V]) Add(key K, value V) {
	l.list.Store(key, value)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (e *list[K, V]) Range(f func(key K, value V) bool) {
	e.list.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (e *list[K, V]) Has(key K) bool {
	ok, _ := e.Load(key)
	return ok
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (l *list[K, V]) Load(key K) (ok bool, value V) {
	if val, ok := l.list.Load(key); ok {
		return true, val.(V)
	}

	return false, value
}

// Remove deletes the value for a key.
func (l *list[K, V]) Remove(key K) {
	l.list.Delete(key)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (e *list[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return e.list.CompareAndDelete(key, old)
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (e *list[K, V]) CompareAndSwap(key K, old, new V) bool {
	return e.list.CompareAndSwap(key, old, new)
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (e *list[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if val, loaded := e.list.LoadAndDelete(key); loaded {
		return val.(V), loaded
	}
	return value, false
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (e *list[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	val, loaded := e.list.LoadOrStore(key, value)
	return val.(V), loaded
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (e *list[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	prev, loaded := e.list.Swap(key, value)
	return prev.(V), loaded
}
