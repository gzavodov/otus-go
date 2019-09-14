package hw4

import (
	"regexp"
	"sort"
	"strings"
)

// Counter is auxiliary structure for counting frequency of words.
type Counter struct {
	Word  string
	Count int
}

// TextAnalyzer is pretty simple text analyzer. Provides only one function GetMostCommonWords.
type TextAnalyzer struct{}

// Regular expression for parsing words in text.
// A word is a sequence of one more letters, numbers, dashes or underscores.
var wordRegexp = regexp.MustCompile(`[\p{L}\p{N}\p{Pd}_]+`)

// GetMostCommonWords returns most frequent words in string defined by 'text' argument.
// Maximum words quantity in result defined by 'limit' argument. If limit <= 0 all the words of text will be returned.
// The result is sorted by frequency of words in the text in descending order.
// If words in result have same frequency they will be sorted by alphabetical ascending order.
func (a *TextAnalyzer) GetMostCommonWords(text string, limit int) []string {
	//Split the source text into words.
	words := wordRegexp.FindAllString(text, -1)

	// Map provides quick access to counter object by word
	counterMap := map[string]*Counter{}
	// List is required for final sorting by frequency of words in the text
	counterList := []*Counter{}

	for _, word := range words {
		// Case of word is ignored. The words 'The' and 'the' are treated as same.
		word = strings.ToLower(word)
		counter, isExist := counterMap[word]

		if isExist {
			counter.Count++
		} else {
			counter = &Counter{Word: word, Count: 1}
			counterList = append(counterList, counter)
			counterMap[word] = counter
		}
	}

	// Sorting by frequency of words is reverted to provide descending order.
	// It uses operation "greater then" instead of "less than".
	// If words have same frequency they will be sorted by natural alphabetical (ascending) order.
	sort.SliceStable(
		counterList,
		func(i, j int) bool {
			if counterList[i].Count == counterList[j].Count {
				return counterList[i].Word < counterList[j].Word
			}
			return counterList[i].Count > counterList[j].Count
		},
	)

	// Limiting the result by a given number of words.
	rightBound := len(counterList)
	if limit > 0 && rightBound > limit {
		rightBound = limit
	}

	result := []string{}
	for i := 0; i < rightBound; i++ {
		result = append(result, counterList[i].Word)
	}
	return result
}
