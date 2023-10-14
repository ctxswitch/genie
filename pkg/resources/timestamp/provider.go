package timestamp

import "time"

// Provider wraps the time interface and is used to provide a time source.
// This is useful for testing.
type Provider interface {
	Now() time.Time
}

// RealTime is the default time provider and is used to provide the current time.
type RealTime struct{}

// Now returns the current time.
// I really want to add some sort option for jitter here as in +- random milliseconds to simulate
// out of order requests across multiple systems.
func (RealTime) Now() time.Time {
	return time.Now()
}

// TestTime is a time provider used for testing.
type TestTime struct {
	Provider,
	override string
}

// NewTestTime returns a new test time provider from a formatted string.
func NewTestTime(ts string) *TestTime {
	return &TestTime{
		override: ts,
	}
}

// Now returns the "current" time that the test time provider is set to.
func (t *TestTime) Now() time.Time {
	now, _ := time.Parse(time.RFC3339Nano, t.override)
	return now
}
