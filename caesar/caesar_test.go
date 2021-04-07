package caesar

import (
	"errors"
	"testing"
)

func TestDecipher(t *testing.T) {
	cases := []struct {
		deciphered string
		ciphered   string
		leaked     []string
		err        error
	}{
		{"abcxyz", "fghcde", []string{"abcxyz"}, nil},
		{"adc, xyz - 1hj#", "uxw, rst - 1bd#", []string{"xyz"}, nil},
		{"adc, xyz - 1hj#", "uxw, rst - 1bd#", []string{"1hj"}, nil},
		{"adc, xyz - 1hj#", "uxw, rst - 1bd#", []string{"1hj"}, nil},
		{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"Yberz vcfhz qbybe fvg nzrg, pbafrpgrghe nqvcvfpvat ryvg, frq qb rvhfzbq grzcbe vapvqvqhag hg ynober rg qbyber zntan nyvdhn.",
			[]string{"consectetur", "elit", "incididunt", "et", "adipiscing"},
			nil,
		},
		{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"Yberz vcfhz qbybe fvg nzrg, pbafrpgrghe nqvcvfpvat ryvg, frq qb rvhfzbq grzcbe vapvqvqhag hg ynober rg qbyber zntan nyvdhn.",
			[]string{"abc", "def"},
			errors.New("key not found for leaked words"),
		},
	}
	for _, c := range cases {
		deciphered, err := Decipher(c.leaked, c.ciphered)
		switch {
		case c.err != nil && err == nil:
			t.Errorf("\nExpected to get error \"%s\"", c.err)
		case c.err == nil && err != nil:
			t.Errorf("\nDid not expect to get error \"%s\"", err)
		case c.err != nil && err.Error() != c.err.Error():
			t.Errorf("\nExpected to get error\n\t\"%s\"\ngot\n\t\"%s\"\n", c.err, err)
		case c.err == nil && deciphered != c.deciphered:
			t.Errorf("\nExpected to decipher\n\t\"%s\"\nwith leaked words %s to\n\t\"%s\"\ngot\n\t\"%s\"\n", c.ciphered, c.leaked, c.deciphered, deciphered)
		}
	}
}

func TestCipher(t *testing.T) {
	cases := []struct {
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
	for _, c := range cases {
		ciphered := Cipher(c.key, c.deciphered)
		if ciphered != c.ciphered {
			t.Errorf("\nExpected to cipher\n\t\"%s\"\nwith key %d to\n\t\"%s\"\ngot\n\t\"%s\"\n", c.deciphered, c.key, c.ciphered, ciphered)
		}
	}
}
