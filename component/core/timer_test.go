package core_test

import (
	"testing"
	"time"

	"github.com/dgf/protzi/component/core"
)

func TestTimer_Run(t *testing.T) {
	durations := make(chan time.Duration)
	stamps := make(chan time.Time)

	timer := &core.Timer{
		Duration: durations,
		Stamp:    stamps,
	}
	go timer.Run()

	durations <- 1 * time.Nanosecond
	select {
	case <-stamps: // OK
	case <-time.After(1 * time.Millisecond):
		t.Error("Timed out!")
	}
}
