package timer

import (
	"fmt"
	"pomodoro/types"
	"strconv"
	"time"
)

type State int

const (
	Stopped State = iota
	Running
	Paused
)

type Timer struct {
	Id        string
	Name      string
	Started   time.Time
	Remaining time.Duration
	State     State
	UpdateCh  chan string
}

var (
	Timers  = make(map[string]*Timer)
	TimerCh = make(chan *Timer)
	// UpdateCh = make(chan string)
)

func NewTimer(duration time.Duration, flags types.Flags) (*Timer, error) {
	timer := &Timer{
		Remaining: duration,
		State:     Stopped,
		// UpdateCh:  make(chan string),
	}
	timer.Id = strconv.Itoa(len(Timers) + 1)
	timer.Name = flags.Name
	timer.UpdateCh = make(chan string)
	Timers[timer.Id] = timer
	timer.Start()
	fmt.Println(timer)
	return timer, nil
}

func (t *Timer) Start() {
	t.State = Running
	t.Started = time.Now()
	go t.countdown()
}

func (t *Timer) Stop() {
	// t.State = Stopped
	t.UpdateCh <- "stop"
}

func (t *Timer) Pause() {
	// t.State = Paused
	// t.Remaining = t.Remaining - time.Since(t.Started)
	t.UpdateCh <- "pause"
}

func (t *Timer) Resume() {
	t.State = Running
}

func (t *Timer) Restart() {
	t.State = Stopped
}

func GetTimers() map[string]*Timer {
	return Timers
}

func GetTimer(id string) *Timer {
	return Timers[id]
}

func (t *Timer) countdown() {
	c := time.Tick(t.Remaining)
	select {
	case <-c:
		t.State = Stopped
		TimerCh <- t
	case msg := <-t.UpdateCh:
		fmt.Println(msg)
		t.Remaining = t.Remaining - time.Since(t.Started)
		if msg == "pause" {
			t.State = Paused
		} else if msg == "stop" {
			t.State = Stopped
		}
	}
	fmt.Println("Timer Stopped")
	return
}
