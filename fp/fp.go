package fp

// A for A in types notation
type A interface{}

// B for B in types notation
type B interface{}

// Predicate function
type Predicate func(el A) bool

// Mapper function
type Mapper func(el A) B

// Reducer function
type Reducer func(acc B, el A) B

// Filter arr according to predicate
func Filter(predicate Predicate, arr []A) []A {
	res := make([]A, 0)
	for _, el := range arr {
		if predicate(el) {
			res = append(res, el)
		}
	}
	return res
}

// FilterC curries Filter
func FilterC(predicate Predicate) func(arr []A) []A {
	return func(arr []A) []A {
		return Filter(predicate, arr)
	}
}

// Map arr according to mapper
func Map(mapper Mapper, arr []A) []B {
	res := make([]B, len(arr))
	for i, el := range arr {
		res[i] = mapper(el)
	}
	return res
}

// MapC curries Map
func MapC(mapper Mapper) func(arr []A) []B {
	return func(arr []A) []B {
		return Map(mapper, arr)
	}
}

// Reduce arr to some value with reducer
func Reduce(reducer Reducer, start B, arr []A) B {
	res := start
	for _, el := range arr {
		res = reducer(res, el)
	}
	return res
}

// ReduceC curries Reduce
func ReduceC(reducer Reducer, start B) func(arr []A) B {
	return func(arr []A) B {
		return Reduce(reducer, start, arr)
	}
}

// ReduceC2 curries ReduceC
func ReduceC2(reducer Reducer) func(start B) func(arr []A) B {
	return func(start B) func(arr []A) B {
		return ReduceC(reducer, start)
	}
}
