package timer

import (
	"fmt"
	"pomodoro/types"
	"time"

	"github.com/google/uuid"
)

type State int

const (
	Stopped State = iota
	Running
	Paused
	Break
)

type Timer struct {
	Id        string
	Name      string
	Duration  time.Duration
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
		Duration:  duration,
		Remaining: duration,
		State:     Stopped,
		// UpdateCh:  make(chan string),
	}
	timer.Id = uuid.New().String()
	timer.Name = flags.Name
	timer.UpdateCh = make(chan string)
	Timers[timer.Id] = timer
	timer.Start()
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
	t.UpdateCh <- "pause"
}

func (t *Timer) Resume() {
	t.Start()
}

func (t *Timer) Restart() {
	t.Remaining = t.Duration
	t.Start()
}

func GetTimers() map[string]*Timer {
	return Timers
}

func GetTimer(id string) *Timer {
	return Timers[id]
}

func DeleteTimer(id string) {
	delete(Timers, id)
}

func (t *Timer) countdown() {
	c := time.Tick(t.Remaining)
	select {
	case <-c:
		// fmt.Println("Timer Stopped")
		t.State = Stopped
		TimerCh <- t
	case msg := <-t.UpdateCh:
		fmt.Println("Message:", msg)
		t.Remaining = t.Remaining - time.Since(t.Started)
		if msg == "pause" {
			t.State = Paused
			t.Remaining = t.Remaining - time.Since(t.Started)
		} else if msg == "stop" {
			t.State = Stopped
		}
	}
	// fmt.Println("Countdown Exit")
}
