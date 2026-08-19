package conv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrCannotConvertToUint64         = errors.New("cannot convert to uint64")
	ErrCannotConvertNegativeToUint64 = errors.New("cannot convert negative to uint64")
)

func ToUint64[T any](d T) (uint64, error) {
	switch v := any(d).(type) {
	case uint64:
		return v, nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf(`int is "%d" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf(`int8 is "%d" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf(`int16 is "%d" and %w`, v, ErrCannotConvertNegativeToUint64)

		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf(`int32 is "%d" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf(`int64 is "%d" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case float32:
		if v < 0 {
			return 0, fmt.Errorf(`float32 is "%f" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf(`float64 is "%f" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(v), nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf(`string is "%s" and %w`, v, ErrCannotConvertNegativeToUint64)
		}
		return uint64(n), nil
	default:
		// Try to handle pointer types
		rv := reflect.ValueOf(d)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			return ToUint64(rv.Elem().Interface())
		}

		return 0, errors.Join(
			fmt.Errorf(`%w "%T"`, ErrUnsupportedType, d),
			ErrCannotConvertToUint64,
		)

	}
}
