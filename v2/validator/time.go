package validator

import (
	"errors"
	"time"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
)

func (v *validator) Time(ctx corectx.Context, t time.Time, field string) error {
	if t.IsZero() {
		return errors.New(v.t.S(ctx, "TimeNotValid", map[string]string{"Field": field}))
	}
	return nil
}

func (v *validator) TimeBetween(ctx corectx.Context, start, end time.Time) error {
	if err := v.Time(ctx, start, "start"); err != nil {
		return err
	}

	if err := v.Time(ctx, end, "end"); err != nil {
		return err
	}

	if end.Before(start) {
		return errors.New(v.t.S(ctx, "EndTimeAfterStart", nil))
	}

	return nil
}
