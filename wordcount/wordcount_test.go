package wordcount

import "testing"

var cases = []struct {
	source string
	result string
}{
	{
		"один, два - это 2, три один два, много слов: один",
		"один(3) два(2) это(1) 2(1) три(1) много(1) слов(1)",
	},
}

func TestCount(t *testing.T) {
	for _, c := range cases {
		result := Count(c.source)
		if result != c.result {
			t.Errorf("\nExpected\n\t\"%s\"\nto result in\n\t\"%s\"\ngot:\n\t\"%s\"\n", c.source, c.result, result)
		}
	}
}
