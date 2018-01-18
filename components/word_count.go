package components

import (
	"strings"
	"text/scanner"
)

// WordCount component cleans text input and counts word occurrences.
type WordCount struct {
	Text   <-chan string
	Counts chan<- map[string]int
}

// cleanup replacements
var replacer = strings.NewReplacer("\f", " ", "\r", " ", "\t", " ")

func (wc *WordCount) Run() {
	for t := range wc.Text {

		// init scanner with cleaned text
		s := &scanner.Scanner{}
		s.Init(strings.NewReader(replacer.Replace(t)))
		s.Mode = scanner.ScanIdents

		// scan and count
		c := map[string]int{}
		for r := s.Scan(); r != scanner.EOF; r = s.Scan() {
			c[s.TokenText()]++
		}

		// write result
		wc.Counts <- c
	}
}
