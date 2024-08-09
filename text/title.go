package text

import "strings"

// Palavras que não devem ser capitalizadas, a menos que sejam a primeira palavra da string
var nonCapitalizedWords = map[string]bool{
	"do":  true,
	"da":  true,
	"de":  true,
	"e":   true,
	"a":   true,
	"ou":  true,
	"no":  true,
	"na":  true,
	"dos": true,
	"das": true,
}

// Title capitaliza a primeira letra de cada palavra, exceto as que estão na lista de exceções
func Title(s string) string {
	if len(s) == 0 {
		return s
	}

	words := strings.Fields(s)
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		first := i == 0
		ignore := !nonCapitalizedWords[lowerWord]
		toUpper := first || ignore

		if toUpper {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		} else {
			words[i] = lowerWord
		}
	}

	return strings.Join(words, " ")
}
