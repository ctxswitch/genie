package timestamp

import (
	"fmt"
	"time"
)

type Timestamp struct {
	// format is a textual representation of the time value formatted
	// according to the layout defined by the argument.  It follows
	// the go representation for time formats.  There are several mapped
	// values that will convert to the go equivalent including: rfc3339,
	// rfc3999nano, unix, and unixnano.
	format    string
	timestamp string
	provider  TimestampProvider
}

func New(cfg Config) *Timestamp {
	return &Timestamp{
		format:    cfg.Format,
		timestamp: cfg.Timestamp,
		provider:  RealTime{},
	}
}

func (t *Timestamp) WithProvider(provider TimestampProvider) *Timestamp {
	t.provider = provider
	return t
}

func (t *Timestamp) Get() string {
	if t.timestamp != "" {
		return t.timestamp
	}

	now := t.provider.Now()
	var format string
	switch t.format {
	case "unix":
		return fmt.Sprint(now.Unix())
	case "unixnano":
		return fmt.Sprint(now.UnixNano())
	case "rfc3339":
		format = time.RFC3339
	case "rfc3339nano":
		format = time.RFC3339Nano
	case "rfc1123":
		format = time.RFC1123
	case "rfc1123z":
		format = time.RFC1123Z
	case "common_log":
		format = CommonLogFormat
	default:
		format = t.format
	}

	formatted := now.Format(format)
	// if the format is not valid, it just returns the format and we'll just return
	// the default format.
	if formatted == format {
		return now.Format("")
	}

	return formatted
}
