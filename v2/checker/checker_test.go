package checker

import (
	"strings"
	"testing"
)

type dep interface{ Do() }

type impl struct{}

func (*impl) Do() {}

type target struct {
	dep  dep
	ptr  *impl
	fn   func()
	list []string
	byID map[string]int
	ch   chan int
	port int
}

func complete() target {
	return target{
		dep:  &impl{},
		ptr:  &impl{},
		fn:   func() {},
		list: []string{},
		byID: map[string]int{},
		ch:   make(chan int),
	}
}

func TestMustComplete(t *testing.T) {
	var typedNil *impl

	tests := []struct {
		name  string
		input func() any
		panic string
	}{
		{"completo", func() any { return complete() }, ""},
		{"ponteiro para struct", func() any { c := complete(); return &c }, ""},
		{"interface nil", func() any { c := complete(); c.dep = nil; return c }, "target.dep"},
		{"ponteiro nil em interface", func() any { c := complete(); c.dep = typedNil; return c }, "target.dep"},
		{"ponteiro nil", func() any { c := complete(); c.ptr = nil; return &c }, "target.ptr"},
		{"func nil", func() any { c := complete(); c.fn = nil; return c }, "target.fn"},
		{"slice nil", func() any { c := complete(); c.list = nil; return c }, "target.list"},
		{"map nil", func() any { c := complete(); c.byID = nil; return c }, "target.byID"},
		{"chan nil", func() any { c := complete(); c.ch = nil; return c }, "target.ch"},
		{"não é struct", func() any { return 1 }, "esperava struct"},
		{"ponteiro nil para struct", func() any { return (*target)(nil) }, "esperava struct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				got, _ := recover().(string)
				if tt.panic == "" && got != "" {
					t.Fatalf("panic inesperado: %s", got)
				}
				if !strings.Contains(got, tt.panic) {
					t.Fatalf("panic = %q, esperava conter %q", got, tt.panic)
				}
			}()

			MustComplete(tt.input())
		})
	}
}
