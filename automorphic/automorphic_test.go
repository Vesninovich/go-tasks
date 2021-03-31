package automorphic

import (
	"testing"
)

func TestAutomorphic(t *testing.T) {
	cases := []struct {
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
	}
	for _, c := range cases {
		a := Automorphic(c.x)
		if a != c.a {
			t.Errorf("Expected %d to be automorphic: %t, got %t", c.x, c.a, a)
		}
	}
}
