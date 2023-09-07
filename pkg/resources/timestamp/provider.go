package timestamp

import "time"

type TimestampProvider interface {
	Now() time.Time
}

type RealTime struct{}

func (RealTime) Now() time.Time {
	return time.Now()
}

type TestTime struct {
	TimestampProvider,
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
