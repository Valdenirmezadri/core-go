// Package checker garante que um struct foi construído por completo.
package checker

import (
	"fmt"
	"reflect"
)

// MustComplete dá panic se s não for um struct, ou ponteiro para struct, com
// todos os campos preenchidos. Serve para rodar em todo New: campo nil é erro
// de programação, então o sistema nem sobe. Usa reflection, só no startup.
func MustComplete(s any) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("checker: esperava struct, recebeu %T", s))
	}

	for i := range v.NumField() {
		if isNil(v.Field(i)) {
			panic(fmt.Sprintf("checker: %s.%s está nil", v.Type(), v.Type().Field(i).Name))
		}
	}
}

// isNil também pega ponteiro nil guardado dentro de interface, que == nil não detecta.
func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface:
		return v.IsNil() || isNil(v.Elem())
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
