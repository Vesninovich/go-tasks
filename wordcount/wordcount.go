package wordcount

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type wordCount struct {
	word  string
	count int
}

type wordCountList []wordCount

func Count(source string) string {
	countsMap := buildWordMap(source)
	counts := make(wordCountList, len(countsMap))
	i := 0
	for word, count := range countsMap {
		counts[i] = wordCount{word, count}
		i++
	}
	sort.Stable(counts)
	var res strings.Builder
	for _, c := range counts {
		res.WriteString(fmt.Sprintf("%s(%d) ", c.word, c.count))
	}
	return res.String()[:res.Len()-1]
}

func buildWordMap(source string) map[string]int {
	countsMap := make(map[string]int)
	var (
		inWord    bool
		wordStart int
	)
	for i, c := range source {
		wordStart, inWord = checkWord(source, countsMap, wordStart, i, c, inWord)
	}
	checkWord(source, countsMap, wordStart, len(source), ' ', inWord)
	return countsMap
}

func checkWord(source string, counts map[string]int, start, i int, c rune, inWord bool) (int, bool) {
	wordChar := unicode.In(c, unicode.Letter, unicode.Digit)
	if !wordChar && inWord {
		word := source[start:i]
		counts[word] += 1
		return start, false
	}
	if wordChar && !inWord {
		return i, true
	}
	return start, inWord
}

func (list wordCountList) Len() int {
	return len(list)
}

func (list wordCountList) Less(i, j int) bool {
	return list[j].count < list[i].count
}

func (list wordCountList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
