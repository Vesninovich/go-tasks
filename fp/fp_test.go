package fp_test

import (
	"fmt"
	"fp"
	"math/rand"
	"strconv"
	"testing"
)

type A struct {
	ID     string
	Name   string
	Age    uint
	Female bool
}

type B string

var cases = prepareCases(1000)

func prepareCases(n uint) (cases []fp.A) {
	cases = make([]fp.A, n)
	for i := range cases {
		cases[i] = A{
			ID:     strconv.Itoa(i),
			Name:   makeName(),
			Age:    uint(rand.Intn(151)),
			Female: rand.Intn(1) == 0,
		}
	}
	return
}

func makeName() (name string) {
	l := 3 + rand.Intn(19)
	for i := 0; i < l; i++ {
		name += string(rune('a' + rand.Intn('z'-'a'+1)))
	}
	return
}

func filter(el fp.A) bool {
	a, ok := el.(A)
	return ok && a.Age >= 18 && !a.Female
}

func BenchmarkFilter(b *testing.B) {
	b.Run("Inline everything", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := make([]A, 0)
			for _, el := range cases {
				a, ok := el.(A)
				if ok && a.Age >= 18 && !a.Female {
					res = append(res, a)
				}
			}
		}
	})
	b.Run("Inline loop, separate predicate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := make([]A, 0)
			for _, el := range cases {
				if filter(el) {
					res = append(res, el.(A))
				}
			}
		}
	})
	b.Run("Filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.Filter(filter, cases)
		}
	})
	b.Run("Curried filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.FilterC(filter)(cases)
		}
	})
}

func mapper(el fp.A) fp.B {
	a, ok := el.(A)
	if !ok {
		panic("wtf")
	}
	title := "Mr"
	if a.Female {
		title = "Ms"
	}
	return fmt.Sprintf("%s %s, %d years old", title, a.Name, a.Age)
}

func BenchmarkMap(b *testing.B) {
	b.Run("Inline everything", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := make([]fp.B, len(cases))
			for i, el := range cases {
				a, ok := el.(A)
				if !ok {
					panic("wtf")
				}
				title := "Mr"
				if a.Female {
					title = "Ms"
				}
				res[i] = fmt.Sprintf("%s %s, %d years old", title, a.Name, a.Age)
			}
		}
	})
	b.Run("Inline loop, separate mapper", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := make([]fp.B, len(cases))
			for i, el := range cases {
				res[i] = mapper(el)
			}
		}
	})
	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.Map(mapper, cases)
		}
	})
	b.Run("Curried map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.MapC(mapper)(cases)
		}
	})
}

func reducer(acc fp.B, el fp.A) fp.B {
	a, ok := el.(A)
	if !ok {
		panic("wtf")
	}
	b, ok := acc.(uint)
	if !ok {
		panic("wtf")
	}
	return b + a.Age
}

func BenchmarkReduce(b *testing.B) {
	b.Run("Inline everything", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			acc := uint(0)
			for _, el := range cases {
				a, ok := el.(A)
				if !ok {
					panic("wtf")
				}
				acc += a.Age
			}
		}
	})
	b.Run("Inline loop, separate reducer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			acc := uint(0)
			for _, el := range cases {
				acc = reducer(acc, el).(uint)
			}
		}
	})
	b.Run("Reduce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.Reduce(reducer, uint(0), cases)
		}
	})
	b.Run("Curried reduce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.ReduceC(reducer, uint(0))(cases)
		}
	})
	b.Run("Curried reduce 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fp.ReduceC2(reducer)(uint(0))(cases)
		}
	})
}
