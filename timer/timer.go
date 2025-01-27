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
	Id       int
	Duration time.Duration
	State    State
}

var Timers = make(map[int]*Timer)

func NewTimer(duration time.Duration) error {
	fmt.Print("New Timer Created")
	timer := &Timer{
		Duration: duration,
		State:    Stopped,
	}
	timer.Id = len(Timers) + 1
	Timers[timer.Id] = timer
	return nil
}

func (t *Timer) Start() {
	t.State = Running
}

func (t *Timer) Stop() {
	t.State = Stopped
}

func (t *Timer) Pause() {
	t.State = Paused
}

func (t *Timer) Resume() {
	t.State = Running
}

func (t *Timer) Restart() {
	t.State = Stopped
}

func GetTimers() map[int]*Timer {
	return Timers
}

func GetTimer(id int) *Timer {
	return Timers[id]
}
