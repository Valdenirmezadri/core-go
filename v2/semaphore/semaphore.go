package semaphore

import "sync"

type Semaphore interface {
	Add(delta int)
	Done()
	Wait()
}

type semaphore struct {
	sem chan bool
	wg  sync.WaitGroup
}

func New(limit uint) Semaphore {
	return &semaphore{
		sem: make(chan bool, limit),
		wg:  sync.WaitGroup{},
	}
}

func (s *semaphore) Add(delta int) {
	s.wg.Add(delta)
	s.sem <- true
}

func (s *semaphore) Done() {
	<-s.sem
	s.wg.Done()
}

func (s *semaphore) Wait() {
	s.wg.Wait()
}
