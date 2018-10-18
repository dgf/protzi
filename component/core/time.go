package core

import (
	"time"
)

// Time component ticks a stamp once.
type Time struct {
	Duration <-chan time.Duration
	Stamp    chan<- time.Time
}

// Run stamps once after duration.
func (t *Time) Run() {
	for duration := range t.Duration {
		t.Stamp <- <-time.NewTimer(duration).C
	}
}
