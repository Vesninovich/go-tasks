`go test -bench ./...`

Benchmarking of filter/map/reduce using full inlining, using separate function with inline loop, uncurried and curried implementations.

Results for Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz:
- Curried implementations perform the same as uncurried
- Filter function is 1.33 times slower than inline loop + separate function, 2+ times slower than fully inlined
- Map function performs the same (maybe mapper function is a bit slow so it evens out a bit)
- Reduce function performs 5 times slower than inline loop + separate function and **34 times slower** than fully inlined

I suspect that significant performance issue is type casting. Moreover, type casting proves to be very inconvenient so until Go 2 and
parametrization, these functions are not only slow, but pretty bad for convenience.
