package log

import "time"

type clock interface {
	now() time.Time
}

type simpleClock struct{}

func (s simpleClock) now() time.Time {
	return time.Now()
}
