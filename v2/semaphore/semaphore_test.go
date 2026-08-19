package semaphore

import (
	"sync"
	"testing"
)

func TestSemaphore_AddAndDone(t *testing.T) {
	limit := uint(2)
	s := New(limit)

	var counter int
	var mu sync.Mutex

	// Function to increment counter and call Done
	increment := func() {
		defer s.Done()
		mu.Lock()
		defer mu.Unlock()
		counter++
	}

	s.Add(1)
	go increment()

	s.Add(1)
	go increment()

	s.Wait()

	if counter != 2 {
		t.Errorf("expected counter to be 2, got %d", counter)
	}

}

func TestSemaphore_Wait(t *testing.T) {
	limit := uint(3)
	s := New(limit)

	var counter int
	var mu sync.Mutex

	// Function to increment counter and call Done
	increment := func() {
		defer s.Done()
		mu.Lock()
		defer mu.Unlock()
		counter++
	}

	for i := 0; i < 50; i++ {
		s.Add(1)
		go increment()
	}

	s.Wait()

	if counter != 50 {
		t.Errorf("expected counter to be 50, got %d", counter)
	}

}

func TestSemaphore_Limit(t *testing.T) {
	limit := uint(2)
	s := New(limit)

	var counter int
	var mu sync.Mutex
	done := make(chan bool)

	// Function to increment counter and call Done
	increment := func() {
		defer func() {
			s.Done()
			done <- true
		}()

		mu.Lock()
		counter++
		mu.Unlock()
	}

	for i := 0; i < 3; i++ {
		s.Add(1)
		go increment()
	}

	// Wait for the first two goroutines to complete
	<-done
	<-done

	if counter != 2 {
		t.Errorf("expected counter to be 2 after first two increments, got %d", counter)
	}

	// Wait for the third goroutine to complete
	<-done

	s.Wait()

	if counter != 3 {
		t.Errorf("expected counter to be 3 after all increments, got %d", counter)
	}

}

func TestSemaphore_ConcurrentAccess(t *testing.T) {
	limit := uint(5)
	s := New(limit)

	var counter int
	var mu sync.Mutex
	numRoutines := 10000
	done := make(chan bool, numRoutines)

	// Function to increment counter and call Done
	increment := func() {
		defer s.Done()
		mu.Lock()
		defer mu.Unlock()
		counter++
		done <- true
	}

	for i := 0; i < numRoutines; i++ {
		s.Add(1)
		go increment()
	}

	for i := 0; i < numRoutines; i++ {
		<-done
	}

	s.Wait()

	if counter != numRoutines {
		t.Errorf("expected counter to be %d, got %d", numRoutines, counter)
	}

}
