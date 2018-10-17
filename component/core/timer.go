package core

import (
	"time"
)

// Timer component ticks a stamp once.
type Timer struct {
	Duration <-chan time.Duration
	Stamp    chan<- time.Time
}

// Run stamps once after duration.
func (t *Timer) Run() {
	for duration := range t.Duration {
		t.Stamp <- <-time.NewTimer(duration).C
	}
}
