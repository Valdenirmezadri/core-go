package htmath

import (
	"fmt"
	"testing"
)

func TestDivide(t *testing.T) {
	type testCase[T Number] struct {
		this     T
		by       T
		expected T
		hasError bool
	}

	testCases := []testCase[float64]{
		{10, 2, 5, false},
		{0, 1, 0, false},
		{10, 0, 0, true},
		{-10, 2, -5, false},
		{10, -2, -5, false},
		{10, 3, 10 / 3.0, false}, // Floating-point division
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v / %v", tc.this, tc.by), func(t *testing.T) {
			result, err := Divide(tc.this, tc.by)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if result != tc.expected {
				t.Errorf("expected result: %v, got: %v", tc.expected, result)
			}
		})
	}

	// Testing with int type
	testCasesInt := []testCase[int]{
		{10, 2, 5, false},
		{0, 1, 0, false},
		{10, 0, 0, true},
		{-10, 2, -5, false},
		{10, -2, -5, false},
		{10, 3, 10 / 3, false}, // Integer division
	}

	for _, tc := range testCasesInt {
		t.Run(fmt.Sprintf("%v / %v", tc.this, tc.by), func(t *testing.T) {
			result, err := Divide(tc.this, tc.by)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if result != tc.expected {
				t.Errorf("expected result: %v, got: %v", tc.expected, result)
			}
		})
	}
}
