package slices

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMap_IntToInt(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	result := Map(input, func(n int) int {
		return n * 2
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMap_IntToString(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	result := Map(input, func(n int) string {
		return fmt.Sprintf("%d", n)
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMap_EmptySlice(t *testing.T) {
	input := []int{}
	expected := []int{}

	result := Map(input, func(n int) int {
		return n * 2
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMap_StringToString(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := []string{"A", "B", "C"}

	result := Map(input, func(s string) string {
		return strings.ToUpper(s)
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
