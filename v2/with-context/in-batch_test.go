package withcontext

import (
	"context"
	"errors"
	"testing"
)

func TestInBatch_ChunksAndCallsOncePerBatch(t *testing.T) {
	var batches [][]int

	err := InBatch(context.Background(), []int{1, 2, 3, 4, 5}, 2, func(batch []int) error {
		batches = append(batches, batch)

		return nil
	})
	if err != nil {
		t.Fatalf("InBatch() error = %v", err)
	}

	if len(batches) != 3 {
		t.Fatalf("want 3 batches, got %d: %v", len(batches), batches)
	}

	// the last batch carries the remainder, not a padded chunk
	if len(batches[2]) != 1 || batches[2][0] != 5 {
		t.Errorf("wrong last batch: %v", batches[2])
	}
}

// a non-positive size would advance by zero and loop forever
func TestInBatch_NonPositiveSizeFallsBackToDefault(t *testing.T) {
	batches := 0

	err := InBatch(context.Background(), make([]int, DefaultBatchSize+1), 0, func([]int) error {
		batches++

		return nil
	})
	if err != nil {
		t.Fatalf("InBatch() error = %v", err)
	}

	if batches != 2 {
		t.Errorf("want 2 batches at the default size, got %d", batches)
	}
}

func TestInBatch_EmptyNeverCalls(t *testing.T) {
	called := false

	if err := InBatch(context.Background(), []int{}, 2, func([]int) error { called = true; return nil }); err != nil {
		t.Fatalf("InBatch() error = %v", err)
	}

	if called {
		t.Error("called fn for empty input")
	}
}

func TestInBatch_StopsOnError(t *testing.T) {
	boom := errors.New("boom")
	calls := 0

	err := InBatch(context.Background(), []int{1, 2, 3, 4, 5, 6}, 2, func([]int) error {
		calls++

		return boom
	})

	if !errors.Is(err, boom) {
		t.Fatalf("want the batch error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("kept going after the error: %d calls", calls)
	}
}

// a cancelled context stops between batches instead of hammering on
func TestInBatch_StopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0

	err := InBatch(ctx, []int{1, 2, 3, 4}, 2, func([]int) error {
		calls++

		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}

	if calls != 0 {
		t.Errorf("ran %d batches on a cancelled context", calls)
	}
}

func TestInBatchAccumulate_JoinsResults(t *testing.T) {
	got, err := InBatchAccumulate(context.Background(), []int{1, 2, 3, 4, 5}, 2, func(batch []int) ([]int, error) {
		return batch, nil
	})
	if err != nil {
		t.Fatalf("InBatchAccumulate() error = %v", err)
	}

	if len(got) != 5 {
		t.Errorf("want 5 joined items, got %d: %v", len(got), got)
	}
}

/*
An error drops what was already accumulated.

Half the rows returned alongside an error read as the whole set to any caller
that only checks the slice.
*/
func TestInBatchAccumulate_ErrorDropsPartial(t *testing.T) {
	boom := errors.New("boom")
	calls := 0

	got, err := InBatchAccumulate(context.Background(), []int{1, 2, 3, 4}, 2, func(batch []int) ([]int, error) {
		calls++
		if calls == 2 {
			return nil, boom
		}

		return batch, nil
	})

	if !errors.Is(err, boom) {
		t.Fatalf("want the batch error, got %v", err)
	}

	if len(got) != 0 {
		t.Errorf("returned a partial result with the error: %v", got)
	}
}
