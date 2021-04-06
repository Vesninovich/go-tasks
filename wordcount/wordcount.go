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

type wordMapRecord struct {
	count int
	order int
}

type wordMap map[string]wordMapRecord

type wordCountList []wordCount

func Count(source string) string {
	countsMap := buildWordMap(source)
	counts := make(wordCountList, len(countsMap))
	for word, r := range countsMap {
		counts[r.order] = wordCount{word, r.count}
	}
	sort.Stable(counts)
	var res strings.Builder
	for _, c := range counts {
		res.WriteString(fmt.Sprintf("%s(%d) ", c.word, c.count))
	}
	return res.String()[:res.Len()-1]
}

func buildWordMap(source string) wordMap {
	countsMap := make(wordMap)
	var (
		inWord    bool
		wordStart int
		wordChar  bool
		word      string
		order     int
		record    wordMapRecord
		exists    bool
	)
	for i, c := range source + " " {
		wordChar = unicode.In(c, unicode.Letter, unicode.Digit)
		if !wordChar && inWord {
			word = source[wordStart:i]
			record, exists = countsMap[word]
			if exists {
				countsMap[word] = wordMapRecord{record.count + 1, record.order}
			} else {
				countsMap[word] = wordMapRecord{1, order}
				order++
			}
			inWord = false
		} else if wordChar && !inWord {
			wordStart, inWord = i, true
		}
	}
	return countsMap
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
