package core

import "time"

// Ticker component ticks a stamp after every looped duration.
type Ticker struct {
	Duration <-chan time.Duration
	Stamps   chan<- time.Time
}

// Run stamps every tick.
func (t *Ticker) Run() {
	ticker := time.NewTicker(<-t.Duration)
	for stamp := range ticker.C {
		t.Stamps <- stamp
	}
}
