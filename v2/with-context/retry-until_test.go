package withcontext

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

var attempts OptFunc

func TestRetryUntil_StopsWhenFnStopsAsking(t *testing.T) {
	calls := 0

	got, err := RetryUntil(context.Background(), func() (int, error) {
		calls++
		if calls < 3 {
			return 0, ErrRetryAgain
		}

		return calls, nil
	})
	if err != nil {
		t.Fatalf("RetryUntil() error = %v", err)
	}

	if got != 3 || calls != 3 {
		t.Errorf("want 3 calls and 3 back, got %d calls and %d", calls, got)
	}
}

// no Attempts means no ceiling: the caller retrying a slug collision must keep
// going until it finds a free one
func TestRetryUntil_WithoutAttemptsKeepsGoing(t *testing.T) {
	calls := 0

	if _, err := RetryUntil(context.Background(), func() (int, error) {
		calls++
		if calls < 50 {
			return 0, ErrRetryAgain
		}

		return calls, nil
	}); err != nil {
		t.Fatalf("RetryUntil() error = %v", err)
	}

	if calls != 50 {
		t.Errorf("stopped early at %d calls", calls)
	}
}

func TestRetryUntil_AttemptsCapsTheTries(t *testing.T) {
	calls := 0

	_, err := RetryUntil(context.Background(), func() (int, error) {
		calls++

		return 0, ErrRetryAgain
	}, attempts.Attempts(3))

	if !errors.Is(err, ErrLimitAttempts) {
		t.Fatalf("want ErrLimitAttempts back, got %v", err)
	}

	if calls != 3 {
		t.Errorf("want 3 tries, got %d", calls)
	}
}

// one attempt means one call: no retry at all, and no off-by-one calling twice
func TestRetryUntil_SingleAttemptCallsOnce(t *testing.T) {
	calls := 0

	if _, err := RetryUntil(context.Background(), func() (int, error) {
		calls++

		return 0, ErrRetryAgain
	}, attempts.Attempts(1)); !errors.Is(err, ErrLimitAttempts) {
		t.Fatalf("want ErrLimitAttempts back, got %v", err)
	}

	if calls != 1 {
		t.Errorf("want a single call, got %d", calls)
	}
}

/*
The ceiling reports itself and keeps the reason.

"Attempt limit reached: 3" alone would not say what kept failing, and that text is
what ends up in front of whoever has to fix it. Wrapping both leaves errors.Is
able to find the ceiling and the cause.
*/
func TestRetryUntil_CeilingKeepsTheReason(t *testing.T) {
	recusa := errors.New("gemini respondeu 503")

	_, err := RetryUntil(context.Background(), func() (int, error) {
		return 0, fmt.Errorf("%w: %w", recusa, ErrRetryAgain)
	}, attempts.Attempts(2))

	if !errors.Is(err, ErrLimitAttempts) {
		t.Errorf("did not report the ceiling: %v", err)
	}

	if !errors.Is(err, recusa) {
		t.Errorf("lost the reason we gave up: %v", err)
	}
}

/*
The wait happens between attempts, not after the last one.

Waiting once the ceiling is reached would only delay an error already decided —
three attempts with one second between them cost two seconds, not three.
*/
func TestRetryUntil_DelayWaitsBetweenAttempts(t *testing.T) {
	calls := 0
	inicio := time.Now()

	_, _ = RetryUntil(context.Background(), func() (int, error) {
		calls++

		return 0, ErrRetryAgain
	}, attempts.Attempts(3), attempts.DelayBetwen(1))

	levou := time.Since(inicio)

	if calls != 3 {
		t.Fatalf("want 3 tries, got %d", calls)
	}

	if levou < 2*time.Second {
		t.Errorf("did not wait between attempts: took %v", levou)
	}

	if levou >= 3*time.Second {
		t.Errorf("waited after the last attempt too: took %v", levou)
	}
}

/*
A cancelled context cuts the wait short.

time.Sleep would hold the goroutine for the whole delay after the caller already
gave up on the result.
*/
func TestRetryUntil_CancelCutsTheWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	inicio := time.Now()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err := RetryUntil(ctx, func() (int, error) {
		return 0, ErrRetryAgain
	}, attempts.DelayBetwen(5))

	levou := time.Since(inicio)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}

	if levou >= time.Second {
		t.Errorf("slept through the cancellation: took %v", levou)
	}
}

// without the option there is no wait at all: the collision retry must stay tight
func TestRetryUntil_NoDelayByDefault(t *testing.T) {
	inicio := time.Now()

	calls := 0

	_, _ = RetryUntil(context.Background(), func() (int, error) {
		calls++

		return 0, ErrRetryAgain
	}, attempts.Attempts(50))

	if levou := time.Since(inicio); levou > 100*time.Millisecond {
		t.Errorf("50 tries took %v with no delay configured", levou)
	}
}

// an error that is not the retry signal comes back on the first try, untouched
func TestRetryUntil_OtherErrorReturnsAtOnce(t *testing.T) {
	boom := errors.New("boom")
	calls := 0

	_, err := RetryUntil(context.Background(), func() (int, error) {
		calls++

		return 0, boom
	}, attempts.Attempts(5))

	if !errors.Is(err, boom) {
		t.Fatalf("want the original error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("retried an error that was not the signal: %d calls", calls)
	}
}

func TestRetryUntil_StopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0

	_, err := RetryUntil(ctx, func() (int, error) {
		calls++

		return 0, ErrRetryAgain
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}

	if calls != 0 {
		t.Errorf("ran %d times on a cancelled context", calls)
	}
}
