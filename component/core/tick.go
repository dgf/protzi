package core

import "time"

// Tick component ticks a stamp after every looped duration.
type Tick struct {
	Duration <-chan time.Duration
	Stamps   chan<- time.Time
}

// Run stamps every tick.
func (t *Tick) Run() {
	ticker := time.NewTicker(<-t.Duration)
	for stamp := range ticker.C {
		t.Stamps <- stamp
	}
}
