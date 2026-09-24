package locks

import "sync/atomic"

type Locker interface {
	// Lock tenta pegar a trava. Devolve alreadyLocked quando outro já está com
	// ela — quem não pegou volta na hora, não espera a vez — e a função que
	// solta, para quem pegou
	Lock() (alreadyLocked bool, releaser func())
}

type handler struct {
	lock atomic.Bool
}

func New() Locker {
	return &handler{}
}

func (u *handler) Lock() (alreadyLocked bool, releaser func()) {
	if !u.lock.CompareAndSwap(false, true) {
		return true, func() {}
	}

	return false, u.Release
}

func (u *handler) Release() {
	u.lock.Store(false)
}
