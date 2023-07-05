package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestSinksHttp(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{`
url: http//localhost
headers:
  - name: Content-Type
    value: application/json
`, true},
		{`
url: http://localhost:3000
headers:
  - name: X-Request-Id
    resource: uuid.request_id
`, true},
		{`
url: http://localhost:3000
headers:
  - name: X-Request-Id
    value: 000000000000000000 
    resource: uuid.request_id
`, false},
	}

	for _, tt := range tests {
		cfg := &Http{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}
