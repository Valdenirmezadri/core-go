package text

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	regSpace             = regexp.MustCompile(`\s+`)
	regAlphanumericOnly  = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	regConsecutiveUnders = regexp.MustCompile(`_+`)
)

// Helper defines text transformation and normalization operations.
type Helper interface {
	// ToFilter normalizes and concatenates strings for use in search filters,
	// removing accents and converting to lowercase.
	ToFilter(list ...string) (term string)
	// TrimAndLower removes leading/trailing whitespace and converts to lowercase.
	TrimAndLower(s string) string
	// Trim removes leading and trailing whitespace.
	Trim(s string) string
	// Capitalize uppercases only the first letter of the string.
	Capitalize(s string) string
	// Title capitalizes each word, except common prepositions and articles.
	Title(s string) string
	// Email normalizes an email address by trimming and lowering.
	Email(s string) string
	// HasSpaces reports whether s contains any whitespace.
	HasSpaces(s string) bool
	// SnakeCase converts s to snake_case, removing accents and special characters.
	SnakeCase(s string) string
	// RemoveAccents strips diacritical marks from s.
	RemoveAccents(s string) string
}

type helper struct{}

// New creates a new Helper.
func New() Helper {
	return &helper{}
}

// ToFilter normalizes and concatenates strings for use in search filters.
func (t helper) ToFilter(list ...string) (term string) {
	for i, item := range list {
		result := t.toFilter(item)

		if i == 0 {
			term = result
			continue
		}

		term += " " + result
	}

	return strings.TrimSpace(term)
}

// toFilter normalizes a single string: removes accents, diacritics and lowercases.
func (t helper) toFilter(s string) string {
	normalized := norm.NFKD.String(s)
	var b strings.Builder
	for _, r := range normalized {
		if !unicode.IsMark(r) {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(strings.ToLower(b.String()))
}

// RemoveAccents strips diacritical marks from s.
func (helper) RemoveAccents(s string) string {
	return RemoveAccents(s)
}

// TrimAndLower removes whitespace and converts to lowercase.
func (t helper) TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Trim removes leading and trailing whitespace.
func (t helper) Trim(s string) string {
	return strings.TrimSpace(s)
}

// Capitalize uppercases only the first letter of the string.
func (t helper) Capitalize(s string) string {
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}
	return s
}

// Title capitalizes each word, except common prepositions and articles.
func (t helper) Title(s string) string {
	return titleCase(t.Trim(s))
}

// Email normalizes an email address by trimming and lowering.
func (t helper) Email(s string) string {
	return t.TrimAndLower(s)
}

// HasSpaces reports whether s contains any whitespace.
func (t helper) HasSpaces(s string) bool {
	return regSpace.MatchString(s)
}

// SnakeCase converts s to snake_case, removing accents and special characters.
func (t helper) SnakeCase(s string) string {
	s = t.RemoveAccents(s)
	s = t.spaceToUnderscore(s)
	s = t.keepAlphanumericAndUnderscore(s)
	s = t.collapseUnderscores(s)
	s = t.TrimAndLower(s)
	return s
}

// spaceToUnderscore replaces whitespace runs with a single underscore.
func (t helper) spaceToUnderscore(s string) string {
	s = regSpace.ReplaceAllString(s, "_")
	s = t.collapseUnderscores(s)
	s = strings.Trim(s, "_")
	return s
}

// collapseUnderscores reduces consecutive underscores to a single one.
func (t helper) collapseUnderscores(s string) string {
	return regConsecutiveUnders.ReplaceAllString(s, "_")
}

// keepAlphanumericAndUnderscore removes all characters except letters, digits and underscores.
func (t helper) keepAlphanumericAndUnderscore(s string) string {
	return regAlphanumericOnly.ReplaceAllString(s, "")
}
