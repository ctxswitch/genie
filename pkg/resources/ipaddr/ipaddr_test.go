package ipaddr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendRange(t *testing.T) {
	expected := []string{
		"10.0.0.1",
		"10.0.0.2",
		"10.0.0.3",
		"10.0.0.4",
		"10.0.0.5",
		"10.0.0.6",
	}

	ips := make([]string, 0)
	cidr := "10.0.0.0/29"
	ips = appendRange(ips, cidr)
	assert.Equal(t, expected, ips)
}
