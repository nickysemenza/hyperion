package clock

import "time"

//Clock is for time
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
}

//RealClock is a normal time.Time
type RealClock struct{}

//Now returns time.Now
func (RealClock) Now() time.Time { return time.Now() }

//Sleep returns time.Sleep
func (RealClock) Sleep(d time.Duration) { time.Sleep(d) }

//After returns time.After
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
