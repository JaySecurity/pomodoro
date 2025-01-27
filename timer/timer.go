package timer

import (
	"fmt"
	"time"
)

type State int

const (
	Stopped State = iota
	Running
	Paused
)

type Timer struct {
	Duration time.Duration
	State    State
}

func NewTimer(duration time.Duration) *Timer {
	fmt.Print("New Timer Created")
	return &Timer{
		Duration: duration,
		State:    Stopped,
	}
}
