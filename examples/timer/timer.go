package main

import (
	"os"
	"strconv"
	"time"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
)

// go run timer.go 7 2
func main() {
	if len(os.Args) < 2 {
		return
	}

	// validate duration seconds
	duration, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// validate interval seconds, defaults to 1
	interval := 1
	if len(os.Args) > 2 {
		i, err := strconv.Atoi(os.Args[2])
		if err == nil {
			interval = i
		}
	}

	// create
	net := protzi.New("timer")
	net.Add("timer", &core.Time{})
	net.Add("ticker", &core.Tick{})
	net.Add("output", &core.Print{})

	// bind
	durations := make(chan time.Duration)
	intervals := make(chan time.Duration)
	out := make(chan bool)
	end := make(chan time.Time)
	net.In("timer.Duration", durations)
	net.In("ticker.Duration", intervals)
	net.Out("output.Printed", out)
	net.Out("timer.Stamp", end)

	// connect and run
	net.Connect("timer.Stamp", "output.Message")
	net.Connect("ticker.Stamps", "output.Message")

	// flow
	durations <- time.Duration(duration) * time.Second
	intervals <- time.Duration(interval) * time.Second

	// discard output response and await the end
	ended := false
	for !ended {
		select {
		case <-out:
			// discard
		case <-end:
			<-out // timer stamp
			ended = true
		}
	}
}
