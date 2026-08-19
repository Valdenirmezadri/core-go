package validator

import (
	"fmt"
	"time"
)

func (validator) Time(time time.Time, field string) error {
	if time.IsZero() {
		return fmt.Errorf("%s time is not valid", field)
	}

	return nil
}

func (v *validator) TimeBetween(start, end time.Time) error {
	if err := v.Time(start, "start"); err != nil {
		return err
	}

	if err := v.Time(end, "end"); err != nil {
		return err
	}

	if end.Before(start) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}
