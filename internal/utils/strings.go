package utils

import (
	"strings"
	"unicode"
)

func MixCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	upper := true
	for _, r := range s {
		if !unicode.IsLetter(r) {
			b.WriteRune(r)
			continue
		}

		if upper {
			b.WriteRune(unicode.ToUpper(r))
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
		upper = !upper
	}

	return b.String()
}
