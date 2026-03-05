package conv

import (
	"errors"
	"testing"
)

func TestToUint64(t *testing.T) {
	type testCase struct {
		input    any
		expected uint64
		wantErr  error
	}

	var (
		negInt           = -5
		negInt8  int8    = -8
		negInt16 int16   = -16
		negInt32 int32   = -32
		negInt64 int64   = -64
		negF32   float32 = -1.5
		negF64   float64 = -2.5
		ptrInt           = 42
		ptrNeg           = -42
	)

	tests := []testCase{
		// uints
		{uint(10), 10, nil},
		{uint8(11), 11, nil},
		{uint16(12), 12, nil},
		{uint32(13), 13, nil},
		{uint64(14), 14, nil},

		// ints positivos
		{int(15), 15, nil},
		{int8(16), 16, nil},
		{int16(17), 17, nil},
		{int32(18), 18, nil},
		{int64(19), 19, nil},

		// ints negativos
		{int(-5), 0, ErrCannotConvertNegativeToUint64},
		{int8(-8), 0, ErrCannotConvertNegativeToUint64},
		{int16(-16), 0, ErrCannotConvertNegativeToUint64},
		{int32(-32), 0, ErrCannotConvertNegativeToUint64},
		{int64(-64), 0, ErrCannotConvertNegativeToUint64},

		// floats positivos
		{float32(20.9), 20, nil},
		{float64(21.9), 21, nil},

		// floats negativos
		{float32(-1.5), 0, ErrCannotConvertNegativeToUint64},
		{float64(-2.5), 0, ErrCannotConvertNegativeToUint64},

		// string válida
		{"123", 123, nil},
		// string inválida
		{"abc", 0, ErrCannotConvertNegativeToUint64},
		// string negativa
		{"-5", 0, ErrCannotConvertNegativeToUint64},

		// ponteiros positivos
		{&ptrInt, 42, nil},
		// ponteiros negativos
		{&ptrNeg, 0, ErrCannotConvertNegativeToUint64},
		{&negInt, 0, ErrCannotConvertNegativeToUint64},
		{&negInt8, 0, ErrCannotConvertNegativeToUint64},
		{&negInt16, 0, ErrCannotConvertNegativeToUint64},
		{&negInt32, 0, ErrCannotConvertNegativeToUint64},
		{&negInt64, 0, ErrCannotConvertNegativeToUint64},
		{&negF32, 0, ErrCannotConvertNegativeToUint64},
		{&negF64, 0, ErrCannotConvertNegativeToUint64},

		// tipo não suportado
		{struct{}{}, 0, ErrUnsupportedType},
	}

	for i, tc := range tests {
		t.Run(string(rune(i)), func(t *testing.T) {
			got, err := ToUint64(tc.input)
			if tc.wantErr != nil {
				if err == nil {
					t.Errorf("expected error, got nil (input: %#v)", tc.input)
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v (input: %#v)", tc.wantErr, err, tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v (input: %#v)", err, tc.input)
				}
				if got != tc.expected {
					t.Errorf("expected %d, got %d (input: %#v)", tc.expected, got, tc.input)
				}
			}
		})
	}
}
