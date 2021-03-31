package automorphic

import (
	"strconv"
)

func Automorphic(n uint64) bool {
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
