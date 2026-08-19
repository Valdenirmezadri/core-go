package locks

import "sync/atomic"

type Locker interface {
	Lock() (holding bool, releaser func())
}

type handler struct {
	lock *uint32
}

func New() Locker {
	var r uint32 = 0
	return &handler{lock: &r}
}

func NewwithRef(ref *uint32) Locker {
	return &handler{lock: ref}
}

func (u *handler) Lock() (alreadyLocked bool, releaser func()) {
	if !atomic.CompareAndSwapUint32(u.lock, 0, 1) {
		return true, func() {}
	}

	return false, u.Release
}

func (u *handler) Release() {
	atomic.StoreUint32(u.lock, 0)
}
