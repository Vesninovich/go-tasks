package compress

import (
	"fmt"
	"strconv"
	"strings"
)

func Compress(source string) string {
	var (
		start      int
		current    rune
		compressed strings.Builder
	)
	compressed.Grow(len(source))
	for i, c := range source {
		if c != current {
			compressed.WriteString(compressedPart(source, current, start, i))
			current, start = c, i
		}
	}
	compressed.WriteString(compressedPart(source, current, start, len(source)))
	return compressed.String()
}

func Decompress(compressed string) string {
	var (
		skip         int
		decompressed strings.Builder
	)
	decompressed.Grow(len(compressed))
	for i, c := range compressed {
		if skip > 0 {
			skip--
			continue
		}
		if c == '#' {
			var (
				char  rune
				count int
			)
			_, err := fmt.Sscanf(compressed[i:], "#%d#%c", &count, &char)
			if err == nil {
				decompressed.WriteString(strings.Repeat(string(char), count))
				skip = len(strconv.Itoa(count)) + 2
				continue
			}
		}
		decompressed.WriteRune(c)
	}
	return decompressed.String()
}

func compressedPart(source string, current rune, start, i int) string {
	count := (i - start) / runeSize(current)
	if count > 4 {
		return fmt.Sprintf("#%d#%c", count, current)
	}
	return source[start:i]
}

func runeSize(r rune) int {
	switch {
	case r < 1<<8:
		return 1
	case r < 1<<16:
		return 2
	case r < 1<<24:
		return 3
	default:
		return 4
	}
}
