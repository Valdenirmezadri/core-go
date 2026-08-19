package text

import "strings"

// nonCapitalizedWords lists prepositions and articles that should remain lowercase
// in title case, unless they are the first word.
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
	"an":  true,
	"on":  true,
	"the": true,
	"to":  true,
}

// Title capitalizes each word, except common prepositions and articles.
func Title(s string) string {
	return titleCase(s)
}

// titleCase capitalizes each word, except common prepositions and articles.
func titleCase(s string) string {
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
