package interrupt

import (
	"os"
	"os/signal"
	"syscall"
)

type Interrupt interface {
	Done() <-chan os.Signal
}

type interrupt struct {
	ch chan os.Signal
}

func New() Interrupt {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return interrupt{ch: ch}
}

func (i interrupt) Done() <-chan os.Signal {
	return i.ch
}
