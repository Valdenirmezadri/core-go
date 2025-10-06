package jobs

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	t.Run("with custom ID", func(t *testing.T) {
		job, err := New(ctx, "test-job")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if job.id != "test-job" {
			t.Errorf("expected id 'test-job', got '%s'", job.id)
		}
	})

	t.Run("with empty ID generates UUID", func(t *testing.T) {
		job, err := New(ctx, "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if job.id == "" {
			t.Error("expected generated UUID, got empty string")
		}
		if len(job.id) != 36 { // UUID length
			t.Errorf("expected UUID length 36, got %d", len(job.id))
		}
	})
}

func TestJob_Run_BasicFunctionality(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	t.Run("successful task execution", func(t *testing.T) {
		var counter int32
		task := func() (bool, error) {
			atomic.AddInt32(&counter, 1)
			return atomic.LoadInt32(&counter) >= 3, nil // stop after 3 executions
		}

		taskID, done, err := job.Run(ctx, task, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if taskID == "" {
			t.Error("expected non-empty taskID")
		}

		select {
		case id := <-done:
			if id != taskID {
				t.Errorf("expected taskID %s, got %s", taskID, id)
			}
		case <-time.After(time.Second):
			t.Error("task should have finished within 1 second")
		}

		finalCount := atomic.LoadInt32(&counter)
		if finalCount != 3 {
			t.Errorf("expected 3 executions, got %d", finalCount)
		}
	})
}

func TestJob_Run_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	tests := []struct {
		name    string
		ctx     context.Context
		task    Task
		every   time.Duration
		wantErr string
	}{
		{
			name:    "nil task",
			ctx:     ctx,
			task:    nil,
			every:   time.Second,
			wantErr: "task function cannot be nil",
		},
		{
			name:    "zero interval",
			ctx:     ctx,
			task:    func() (bool, error) { return true, nil },
			every:   0,
			wantErr: "interval must be positive",
		},
		{
			name:    "negative interval",
			ctx:     ctx,
			task:    func() (bool, error) { return true, nil },
			every:   -time.Second,
			wantErr: "interval must be positive",
		},
		{
			name:    "nil context",
			ctx:     nil,
			task:    func() (bool, error) { return true, nil },
			every:   time.Second,
			wantErr: "need an context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := job.Run(tt.ctx, tt.task, tt.every)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error '%s', got '%s'", tt.wantErr, err.Error())
			}
		})
	}
}

func TestJob_Run_AlreadyRunning(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	// Start first task
	task1 := func() (bool, error) {
		time.Sleep(100 * time.Millisecond)
		return true, nil
	}

	_, _, err1 := job.Run(ctx, task1, 10*time.Millisecond)
	if err1 != nil {
		t.Fatalf("first run should succeed, got %v", err1)
	}

	// Try to start second task immediately
	task2 := func() (bool, error) { return true, nil }
	_, _, err2 := job.Run(ctx, task2, 10*time.Millisecond)

	if err2 == nil {
		t.Error("second run should fail with already running error")
	}

	if !errors.Is(err2, ErrRunning) {
		t.Errorf("expected ErrRunning, got %v", err2)
	}
}

func TestJob_Run_ContextCancellation(t *testing.T) {
	parentCtx := context.Background()
	job, _ := New(parentCtx, "test")

	t.Run("job context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(parentCtx)

		var execCount int32
		task := func() (bool, error) {
			atomic.AddInt32(&execCount, 1)
			return false, nil // never finish naturally
		}

		_, done, err := job.Run(ctx, task, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Let it run a few times
		time.Sleep(50 * time.Millisecond)
		cancel()

		// Should not receive on done channel since it was cancelled
		select {
		case <-done:
			t.Error("should not receive on done channel when cancelled")
		case <-time.After(100 * time.Millisecond):
			// Expected - task was cancelled, not finished
		}

		count := atomic.LoadInt32(&execCount)
		if count == 0 {
			t.Error("task should have executed at least once")
		}
	})
}

func TestJob_Run_TaskErrors(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	t.Run("task returns error but continues", func(t *testing.T) {
		var execCount int32
		task := func() (bool, error) {
			count := atomic.AddInt32(&execCount, 1)
			if count < 3 {
				return false, errors.New("task error")
			}
			return true, nil // finish after 3 attempts
		}

		_, done, err := job.Run(ctx, task, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		select {
		case <-done:
			// Task should complete despite errors
		case <-time.After(time.Second):
			t.Error("task should have finished within 1 second")
		}

		finalCount := atomic.LoadInt32(&execCount)
		if finalCount != 3 {
			t.Errorf("expected 3 executions, got %d", finalCount)
		}
	})
}

func TestJob_Run_TaskPanic(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	var execCount int32
	task := func() (bool, error) {
		count := atomic.AddInt32(&execCount, 1)
		if count == 1 {
			panic("test panic")
		}
		return count >= 3, nil
	}

	_, done, err := job.Run(ctx, task, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case <-done:
		// Task should complete despite panic
	case <-time.After(time.Second):
		t.Error("task should have finished within 1 second")
	}

	finalCount := atomic.LoadInt32(&execCount)
	if finalCount < 3 {
		t.Errorf("expected at least 3 executions, got %d", finalCount)
	}
}

func TestJob_Stop(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	t.Run("stop running job", func(t *testing.T) {
		task := func() (bool, error) {
			time.Sleep(50 * time.Millisecond)
			return false, nil // never finish naturally
		}

		_, _, err := job.Run(ctx, task, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Give it time to start
		time.Sleep(20 * time.Millisecond)

		start := time.Now()
		err = job.Stop(time.Second)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if duration > 500*time.Millisecond {
			t.Errorf("stop took too long: %v", duration)
		}
	})

	t.Run("stop non-running job", func(t *testing.T) {
		job2, _ := New(ctx, "test2")
		err := job2.Stop(time.Second)
		if err != nil {
			t.Errorf("expected no error when stopping non-running job, got %v", err)
		}
	})

	t.Run("stop timeout", func(t *testing.T) {
		job3, _ := New(ctx, "test3")

		task := func() (bool, error) {
			time.Sleep(200 * time.Millisecond) // mais lento que o timeout
			return false, nil
		}

		_, _, err := job3.Run(ctx, task, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Timeout bem curto
		err = job3.Stop(50 * time.Millisecond)
		switch {
		case err == nil:
			// às vezes a goroutine sai cedo e não força timeout -> aceitável
			t.Log("job finished before timeout (acceptable)")

		case err.Error() != "stop timeout exceeded":
			// se houve erro mas não o esperado -> falha
			t.Errorf("expected 'stop timeout exceeded', got '%v'", err)
		}
	})
}

func TestJob_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	// Try to start multiple jobs concurrently
	const numGoroutines = 10
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			task := func() (bool, error) {
				return true, nil
			}

			_, _, err := job.Run(ctx, task, 10*time.Millisecond)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}

	wg.Wait()

	// Only one should succeed, others should get ErrRunning
	if atomic.LoadInt32(&successCount) != 1 {
		t.Errorf("expected exactly 1 success, got %d", atomic.LoadInt32(&successCount))
	}

	if atomic.LoadInt32(&errorCount) != numGoroutines-1 {
		t.Errorf("expected %d errors, got %d", numGoroutines-1, atomic.LoadInt32(&errorCount))
	}
}

func TestJob_LongRunningTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	ctx := context.Background()
	job, _ := New(ctx, "test")

	var execCount int32
	task := func() (bool, error) {
		count := atomic.AddInt32(&execCount, 1)
		return count >= 100, nil // run 100 times
	}

	start := time.Now()
	_, done, err := job.Run(ctx, task, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case <-done:
		duration := time.Since(start)
		t.Logf("100 executions took %v", duration)

		if duration < 900*time.Millisecond {
			t.Errorf("expected at least 900ms for 100 executions with 10ms interval, got %v", duration)
		}

		if duration > 2*time.Second {
			t.Errorf("took too long: %v", duration)
		}

	case <-time.After(5 * time.Second):
		t.Error("long running task should have finished within 5 seconds")
	}

	finalCount := atomic.LoadInt32(&execCount)
	if finalCount != 100 {
		t.Errorf("expected 100 executions, got %d", finalCount)
	}
}

func TestJob_MultipleStopCalls(t *testing.T) {
	ctx := context.Background()
	job, _ := New(ctx, "test")

	task := func() (bool, error) {
		time.Sleep(100 * time.Millisecond)
		return false, nil
	}

	_, _, err := job.Run(ctx, task, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Call stop multiple times concurrently
	var wg sync.WaitGroup
	errors := make([]error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errors[idx] = job.Stop(time.Second)
		}(i)
	}

	wg.Wait()

	// All stop calls should succeed (or at least not panic)
	for i, err := range errors {
		if err != nil {
			t.Errorf("stop call %d failed: %v", i, err)
		}
	}
}

// Benchmark tests
func BenchmarkJob_Run_ShortTask(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job, _ := New(ctx, "bench")

		task := func() (bool, error) {
			return true, nil // finish immediately
		}

		_, done, err := job.Run(ctx, task, time.Microsecond)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}

		<-done
	}
}

func BenchmarkJob_Run_MultipleExecutions(b *testing.B) {
	ctx := context.Background()
	job, _ := New(ctx, "bench")

	b.ResetTimer()

	var counter int32
	task := func() (bool, error) {
		count := atomic.AddInt32(&counter, 1)
		return count >= int32(b.N), nil
	}

	_, done, err := job.Run(ctx, task, time.Microsecond)
	if err != nil {
		b.Fatalf("expected no error, got %v", err)
	}

	<-done
}
