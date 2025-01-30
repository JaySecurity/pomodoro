package types

import "time"

type Flags struct {
	Index    string
	Duration string
	Name     string
}

type Timer struct {
	Id        string
	Name      string
	Started   time.Time
	Remaining time.Duration
	State     int
}

var Status = map[int]string{
	0: "Stopped",
	1: "Running",
	2: "Paused",
}
