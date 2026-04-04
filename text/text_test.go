package text

import "testing"

func Test_helper_Trim(t *testing.T) {
	h := New()
	tests := []struct {
		name string
		input string
		want string
	}{
		{name: "no spaces", input: "Valdenir Luíz Mezadri Junior", want: "Valdenir Luíz Mezadri Junior"},
		{name: "both sides", input: " Valdenir Luíz Mezadri Junior ", want: "Valdenir Luíz Mezadri Junior"},
		{name: "left space", input: " Valdenir Luíz Mezadri Junior", want: "Valdenir Luíz Mezadri Junior"},
		{name: "right space", input: "Valdenir Luíz Mezadri Junior ", want: "Valdenir Luíz Mezadri Junior"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.Trim(tt.input); got != tt.want {
				t.Errorf("helper.Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_Title(t *testing.T) {
	h := New()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "leading space", input: " valdenir luíz mezadri junior", want: "Valdenir Luíz Mezadri Junior"},
		{name: "the preposition", input: "Welcome to the Dollhouse ", want: "Welcome to the Dollhouse"},
		{name: "mixed case", input: "a niGht to reMember", want: "A Night to Remember"},
		{name: "dos preposition", input: "dieGo dos santos ", want: "Diego dos Santos"},
		{name: "multiple prepositions", input: "DieGo das dores dos PraZeres ", want: "Diego das Dores dos Prazeres"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.Title(tt.input); got != tt.want {
				t.Errorf("helper.Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_Email(t *testing.T) {
	h := New()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "trailing space", input: "junior@hardtec.srv.br ", want: "junior@hardtec.srv.br"},
		{name: "leading space", input: " junior@hardtec.srv.br", want: "junior@hardtec.srv.br"},
		{name: "both spaces", input: " junior@hardtec.srv.br ", want: "junior@hardtec.srv.br"},
		{name: "capitalized", input: "Junior@hardtec.srv.br ", want: "junior@hardtec.srv.br"},
		{name: "mixed case", input: "JUnior@hardTEC.srv.br ", want: "junior@hardtec.srv.br"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.Email(tt.input); got != tt.want {
				t.Errorf("helper.Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	got := New()
	if got == nil {
		t.Errorf("New() = nil, want not nil")
	}
	if _, ok := got.(*helper); !ok {
		t.Errorf("New() type = %T, want *helper", got)
	}
}

func Test_helper_SnakeCase(t *testing.T) {
	h := New()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "mixed with accents", input: "Olá, mundo! 123 - teste", want: "ola_mundo_123_teste"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.SnakeCase(tt.input); got != tt.want {
				t.Errorf("helper.SnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
