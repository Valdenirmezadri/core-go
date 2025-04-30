package slices

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	// Caso de teste 1: Filtrar números pares
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedEvens := []int{2, 4, 6, 8, 10}
	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	if !reflect.DeepEqual(evens, expectedEvens) {
		t.Errorf("Esperado %v, mas obteve %v", expectedEvens, evens)
	}

	// Caso de teste 2: Filtrar strings que começam com 'A'
	strings := []string{"Apple", "Banana", "Avocado", "Cherry", "Apricot"}
	expectedStartsWithA := []string{"Apple", "Avocado", "Apricot"}
	startsWithA := Filter(strings, func(s string) bool {
		return len(s) > 0 && s[0] == 'A'
	})
	if !reflect.DeepEqual(startsWithA, expectedStartsWithA) {
		t.Errorf("Esperado %v, mas obteve %v", expectedStartsWithA, startsWithA)
	}

	// Caso de teste 3: Filtrar elementos de um slice vazio
	emptyInitialized := []int{}
	expectedEmpty := []int{}
	filteredEmptyInitialized := Filter(emptyInitialized, func(n int) bool {
		return n%2 == 0
	})
	if !reflect.DeepEqual(filteredEmptyInitialized, expectedEmpty) {
		t.Errorf("Esperado %v, mas obteve %v", expectedEmpty, filteredEmptyInitialized)
	}

	// Caso de teste 4: Nenhum elemento satisfaz o predicado
	numbersNone := []int{1, 3, 5, 7, 9}
	expectedNone := []int{}
	filteredNone := Filter(numbersNone, func(n int) bool {
		return n%2 == 0
	})
	if !reflect.DeepEqual(filteredNone, expectedNone) {
		t.Errorf("Esperado %v, mas obteve %v", expectedNone, filteredNone)
	}

	// Caso de teste 5: Todos os elementos satisfazem o predicado
	allNumbers := []int{2, 4, 6, 8, 10}
	expectedAll := []int{2, 4, 6, 8, 10}
	filteredAll := Filter(allNumbers, func(n int) bool {
		return n%2 == 0
	})
	if !reflect.DeepEqual(filteredAll, expectedAll) {
		t.Errorf("Esperado %v, mas obteve %v", expectedAll, filteredAll)
	}
}
