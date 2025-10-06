package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Valdenirmezadri/core-go/safe"
	"github.com/google/uuid"
)

// Task represents a function that can be executed periodically.
// It returns finished=true to stop the job, or finished=false to continue.
type Task func() (finished bool, err error)

var (
	ErrRunning = errors.New("already running")
)

// Job represents a periodic task executor.
type Job struct {
	ctx      context.Context
	id       string
	_running atomic.Bool
	_cancel  safe.Item[context.CancelFunc]
}

// New creates a new Job instance.
func New(ctx context.Context, id string) (*Job, error) {
	if id == "" {
		id = uuid.NewString()
	}

	return &Job{
		ctx:     ctx,
		id:      id,
		_cancel: safe.NewItem[context.CancelFunc](),
	}, nil
}

func (r *Job) Stop(timeout time.Duration) error {
	if !r._running.Load() {
		return nil
	}

	r.cancel()

	done := make(chan struct{})
	go func() {
		for r._running.Load() {
			time.Sleep(10 * time.Millisecond)
		}

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("stop timeout exceeded")
	}
}

func (r *Job) cancel() {
	r._cancel.Update(func(cf context.CancelFunc) context.CancelFunc {
		if cf != nil {
			cf()
		}

		return nil
	})
}

// Run starts executing the task periodically until it finishes or is stopped.
func (r *Job) Run(ctx context.Context, task Task, every time.Duration) (taskID string, done <-chan string, err error) {
	if task == nil {
		return "", nil, errors.New("task function cannot be nil")
	}

	if every <= 0 {
		return "", nil, errors.New("interval must be positive")
	}

	if ctx == nil {
		return "", nil, errors.New("need an context")
	}

	if !r._running.CompareAndSwap(false, true) {
		return "", nil, fmt.Errorf(`job "%s": %w`, r.id, ErrRunning)
	}

	ctx, cancel := context.WithCancel(ctx)
	r._cancel.Set(cancel)
	taskID = uuid.NewString()
	finished := make(chan string, 1)

	go r.run(ctx, taskID, task, finished, every)

	return taskID, finished, nil
}

func (r *Job) run(ctx context.Context, taskID string, task Task, done chan<- string, every time.Duration) {
	defer func() {
		r._running.CompareAndSwap(true, false)
	}()

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			finished, err := r.safeExecute(task)
			if err != nil {
				fmt.Printf("task %s of job %s execution error: %v\n", taskID, r.id, err)
			}

			if finished {
				done <- taskID
				return
			}
		}
	}
}

func (r *Job) safeExecute(task Task) (finished bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			// Log do panic sem quebrar o worker
			fmt.Printf("task of job %s panic recovered: %v\n", r.id, recovered)
		}
	}()

	return task()
}
