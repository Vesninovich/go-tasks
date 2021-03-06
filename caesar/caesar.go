package caesar

import (
	"errors"
	"strings"
)

type empty struct{}

const latLetterAmount = 'Z' - 'A' + 1

// Cipher shifts all latin characters right by `key`,
// as described by Caesar Cipher algorithm
func Cipher(key int, source string) string {
	return strings.Map(rotCharBy(normalizeKey(key)), source)
}

// Decipher tries to guess cipher key from given deciphered words of ciphered text
// and then decipher it by ciphering text by guessed key.
func Decipher(leaked []string, ciphered string) (string, error) {
	leakedWords := getLeakedWords(leaked)
	key, err := getCipherKey(leakedWords, ciphered)
	if err != nil {
		return "", err
	}
	return Cipher(key, ciphered), nil
}

func getLeakedWords(leaked []string) (leakedWords map[string]empty) {
	leakedWords = make(map[string]empty)
	for _, w := range leaked {
		leakedWords[w] = empty{}
	}
	return leakedWords
}

func getCipherKey(leaked map[string]empty, ciphered string) (key int, err error) {
	var (
		inWord    bool
		wordStart int
		wordChar  bool
		word      string
		found     bool
	)
	for i, c := range ciphered + " " {
		wordChar = isASCIIAlpha(c)
		if !wordChar && inWord {
			word = ciphered[wordStart:i]
			key, found = keyFromWord(leaked, word)
			if found {
				return key, nil
			}
			inWord = false
		} else if wordChar && !inWord {
			wordStart, inWord = i, true
		}
	}
	return 0, errors.New("key not found for leaked words")
}

func keyFromWord(leaked map[string]empty, word string) (key int, found bool) {
	var rot string
	for key = 1; key < latLetterAmount; key++ {
		rot = strings.Map(rotCharBy(key), word)
		_, found = leaked[rot]
		if found {
			return
		}
	}
	return
}

func rotCharBy(n int) func(rune) rune {
	return func(r rune) rune {
		if !isASCIILetter(r) {
			return r
		}
		res := r + rune(n)
		if r <= 'Z' && res > 'Z' || r <= 'z' && res > 'z' {
			return res - latLetterAmount
		}
		return res
	}
}

func isASCIILetter(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

func isASCIIAlpha(r rune) bool {
	return isASCIILetter(r) || r >= '0' && r <= '9'
}

func normalizeKey(key int) int {
	key -= latLetterAmount * int(key/latLetterAmount)
	if key < 0 {
		key += latLetterAmount
	}
	return key
}
