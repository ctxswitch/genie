package timestamp

import (
	"fmt"
	"time"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/timestamp"
)

type Timestamp struct {
	format    string
	timestamp string
	provider  timestamp.TimestampProvider
}

func FromConfig(options config.Timestamp) *Timestamp {
	return &Timestamp{
		format:    options.Format,
		timestamp: options.Timestamp,
		provider:  timestamp.RealTime{},
	}
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
		return time.RFC3339
	case "rfc3339nano":
		return time.RFC3339Nano
	case "rfc1123":
		return time.RFC1123
	case "rfc1123z":
		return time.RFC1123Z
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

func (t *Timestamp) WithProvider(provider timestamp.TimestampProvider) *Timestamp {
	t.provider = provider
	return t
}
