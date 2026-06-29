package scheduler

import "time"

// Clock is an indirection over time.Now. Tests use a FakeClock.
type Clock interface {
	Now() time.Time
}

// RealClock is the production Clock.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }