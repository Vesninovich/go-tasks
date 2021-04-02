package caesar

import (
	"strings"
)

const latLetterAmount = 'Z' - 'A' + 1

func Cipher(key int, source string) string {
	return strings.Map(rotCharBy(normalizeKey(key)), source)
}

func Decipher(key int, ciphered string) string {
	return Cipher(latLetterAmount-normalizeKey(key), ciphered)
}

func rotCharBy(n int) func(rune) rune {
	return func(r rune) rune {
		if !isAsciiLetter(r) {
			return r
		}
		r += rune(n)
		if !isAsciiLetter(r) {
			r -= latLetterAmount
		}
		return r
	}
}

func isAsciiLetter(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

func normalizeKey(key int) int {
	key -= latLetterAmount * int(key/latLetterAmount)
	if key < 0 {
		key += latLetterAmount
	}
	return key
}
