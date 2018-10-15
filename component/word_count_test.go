package component_test

import (
	"fmt"
	"sort"

	"github.com/dgf/protzi/component"
)

func ExampleWordCount() {
	text := make(chan string)
	counts := make(chan map[string]int)

	// create and run counter
	go (&component.WordCount{Text: text, Counts: counts}).Run()

	// send text
	text <- "\f\n\r\ttwo\fone\ntwo\r\t"

	// read word counts
	countsByWord := <-counts

	// stringify and sort word counts (needed for output assertion)
	var wordCounts []string
	for word := range countsByWord {
		wordCounts = append(wordCounts, fmt.Sprintf("%s: %d", word, countsByWord[word]))
	}
	sort.Strings(wordCounts)

	// print word counts
	fmt.Println(wordCounts)

	// Output: [one: 1 two: 2]
}
