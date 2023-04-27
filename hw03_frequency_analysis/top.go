package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type frequency struct {
	word  string
	count int
}

func Top10(input string) []string {
	if len(input) == 0 {
		return []string{}
	}

	words := strings.Fields(input)
	wordsCounts := make(map[string]int)

	for _, word := range words {
		wordsCounts[word]++
	}

	frequencies := make([]frequency, 0, len(wordsCounts))
	for key, val := range wordsCounts {
		frequencies = append(frequencies, frequency{word: key, count: val})
	}
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	topCount := 10
	if len(frequencies) < 10 {
		topCount = len(frequencies)
	}

	result := make([]string, topCount)
	for i := 0; i < topCount; i++ {
		result[i] = frequencies[i].word
	}
	return result
}
