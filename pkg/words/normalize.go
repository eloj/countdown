package words

import (
	"strings"

	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeLatin1 lowercases, replaces codepoints beyond 255 with a space.
func NormalizeLatin1(word string) string {
	word = strings.ToLower(word)

	runeMapping := runes.Map(func(r rune) rune {
		if !utf8.ValidRune(r) || unicode.IsControl(r) {
			return ' '
		}

		// We want up to extended ASCII only.
		if r > 255 {
			return ' '
		}

		return r
	})

	tc := transform.Chain(norm.NFKC, runeMapping)
	res, _, err := transform.String(tc, word)
	if err != nil {
		return ""
	}
	return res
}
