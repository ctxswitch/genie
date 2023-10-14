package timestamp

import "time"

type Provider interface {
	Now() time.Time
}

type RealTime struct{}

// I really want to add some sort option for jitter here as in +- random milliseconds to simulate
// out of order requests across multiple systems.
func (RealTime) Now() time.Time {
	return time.Now()
}

type TestTime struct {
	Provider,
	override string
}

func NewTestTime(ts string) *TestTime {
	return &TestTime{
		override: ts,
	}
}

func (t *TestTime) Now() time.Time {
	now, _ := time.Parse(time.RFC3339Nano, t.override)
	return now
}
