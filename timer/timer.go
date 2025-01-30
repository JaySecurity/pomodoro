package timer

import (
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
}

var (
	Timers   = make(map[string]*Timer)
	TimerCh  = make(chan *Timer)
	UpdateCh = make(chan string)
)

func NewTimer(duration time.Duration, flags types.Flags) (*Timer, error) {
	timer := &Timer{
		Remaining: duration,
		State:     Stopped,
	}
	timer.Id = strconv.Itoa(len(Timers) + 1)
	timer.Name = flags.Name
	Timers[timer.Id] = timer
	timer.Start()
	return timer, nil
}

func (t *Timer) Start() {
	t.State = Running
	t.Started = time.Now()
	go countdown(UpdateCh, t)
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

func GetTimers() map[string]*Timer {
	return Timers
}

func GetTimer(id string) *Timer {
	return Timers[id]
}

func runTimer() {
	time.Tick(5 * time.Second)
}

func countdown(ch chan string, t *Timer) {
	c := time.Tick(t.Remaining)
	select {
	case <-c:
		t.State = Stopped
		TimerCh <- t
	case msg := <-ch:
		if msg == "pause" {
			t.State = Paused
			t.Remaining = t.Remaining - time.Since(t.Started)
			break
		} else if msg == "stop" {
			t.State = Stopped
			defer close(ch)
			break
		}
	}
}
