package automorphic

import (
	"testing"
)

var cases = []struct {
	x uint64
	a bool
}{
	{0, true},
	{1, true},
	{3, false},
	{5, true},
	{16, false},
	{25, true},
	{100, false},
	{376, true},
	{90625, true},
	{100000, false},
}

func TestAutomorphicMods(t *testing.T) {
	runTest(t, automorphicMath, "Math")
	runTest(t, automorphicStrings, "Strings")
}

func runTest(t *testing.T, automorphicFunc func(n uint64) bool, testName string) {
	for _, c := range cases {
		a := automorphicFunc(c.x)
		if a != c.a {
			t.Errorf("Expected %d to be automorphic: %t, got %t in test %s\n", c.x, c.a, a, testName)
		}
	}
}
