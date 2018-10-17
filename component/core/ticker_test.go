package core_test

import (
	"testing"
	"time"

	"github.com/dgf/protzi/component/core"
)

func TestTicker_Run(t *testing.T) {
	durations := make(chan time.Duration)
	stamps := make(chan time.Time)
	ticker := &core.Ticker{
		Duration: durations,
		Stamps:   stamps,
	}
	go ticker.Run()

	count := 0
	durations <- 1 * time.Nanosecond
	for count < 2 {
		select {
		case <-stamps:
			count++
		case <-time.After(1 * time.Millisecond):
			t.Error("Timed out!")
		}
	}
}
