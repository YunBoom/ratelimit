package ratelimit

import (
	"time"
)

type sleepClock struct{}

func (s sleepClock) Now() time.Time {
	return time.Now()
}

func (s sleepClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
