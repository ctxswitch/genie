package timestamp

import "time"

const (
	// DefaultTimeFormat is the default time format.
	DefaultTimeFormat = time.RFC3339Nano
	// DefaultTimeTimestamp is the default time timestamp.
	DefaultTimeTimestamp = ""
	// CommonLogFormat is the common log format used by web servers.
	CommonLogFormat = "02/Jan/2006:15:04:05 -0700"
)
