package automorphic

import (
	"math"
	"strconv"
)

// Two different implementations, because why not
func Automorphic(n uint64) bool {
	return automorphicMath(n)
}

func automorphicMath(n uint64) bool {
	if n == 0 {
		return true
	}
	sq := n * n
	pow := int(math.Floor(math.Log10(float64(n)))) + 1
	return sq%uint64(math.Pow10(pow)) == n
}

func automorphicStrings(n uint64) bool {
	sq := n * n
	nString, sqString := strconv.FormatUint(n, 10), strconv.FormatUint(sq, 10)
	nSize, sqSize := len(nString), len(sqString)
	i := nSize - 1
	for ; i >= 0; i-- {
		sqIndex := sqSize - (nSize - i)
		if nString[i] != sqString[sqIndex] {
			return false
		}
	}
	return i < 0
}
