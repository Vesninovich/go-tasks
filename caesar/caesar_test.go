package caesar

import "testing"

var cases = []struct {
	deciphered string
	ciphered   string
	key        int
}{
	{"abcxyz", "fghcde", 5},
	{"abc, xyz - 1hj#", "fgh, cde - 1mo#", 5},
	{"abc, xyz - 1hj#", "uvw, rst - 1bd#", 20},
	{"abc, xyz - 1hj#", "fgh, cde - 1mo#", 265},
	{"abc, xyz - 1hj#", "fgh, cde - 1mo#", -541},
}

func TestCipher(t *testing.T) {
	for _, c := range cases {
		ciphered := Cipher(c.key, c.deciphered)
		if ciphered != c.ciphered {
			t.Errorf("\nExpected to cipher\n\t\"%s\"\nwith key %d to\n\t\"%s\"\ngot\n\t\"%s\"\n", c.deciphered, c.key, c.ciphered, ciphered)
		}
		deciphered := Decipher(c.key, c.ciphered)
		if deciphered != c.deciphered {
			t.Errorf("\nExpected to decipher\n\t\"%s\"\nwith key %d to\n\t\"%s\"\ngot\n\t\"%s\"\n", c.ciphered, c.key, c.deciphered, deciphered)
		}
	}
}
