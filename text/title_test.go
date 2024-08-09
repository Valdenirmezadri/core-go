package text

import (
	"testing"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"rio do sul", "Rio do Sul"},
		{"festa na cidade", "Festa na Cidade"},
		{"o sol e a lua", "O Sol e a Lua"},
		{"rio de janeiro", "Rio de Janeiro"},
		{"", ""},
		{"uma palavra", "Uma Palavra"},
		{"do sol", "Do Sol"},
		{"e", "E"},
		{"DA Lua", "Da Lua"},
		{"de MARTE", "De Marte"},
		{"um teste do sistema", "Um Teste do Sistema"},
		{"golang e divertido", "Golang e Divertido"},
	}

	for _, test := range tests {
		result := Title(test.input)
		if result != test.expected {
			t.Errorf("Title(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
