package compress

import "testing"

var cases = []struct {
	d string
	c string
}{
	{"", ""},
	{"кооооооордината", "к#7#ордината"},
	{"asdfgh", "asdfgh"},
	{"qweeeertyyyyyy", "qweeeert#6#y"},
	{"999999999", "#9#9"},
	{"899999999999998", "8#13#98"},
	{"##a#######b", "##a#7##b"},
}

func TestCompress(t *testing.T) {
	for _, c := range cases {
		comp := Compress(c.d)
		if comp != c.c {
			t.Errorf("Expected compressed %s to be %s, got %s\n", c.d, c.c, comp)
		}
	}
}

func TestDecompress(t *testing.T) {
	for _, c := range cases {
		decomp := Decompress(c.c)
		if decomp != c.d {
			t.Errorf("Expected decompressed %s to be %s, got %s\n", c.c, c.d, decomp)
		}
	}
}
