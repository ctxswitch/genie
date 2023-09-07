package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAll(t *testing.T) {
	path, err := os.Getwd()
	require.NoError(t, err)

	dir := fmt.Sprintf("%s/../../genie.d", path)
	_, err = Load(&LoadOptions{
		Paths:  []string{dir},
		Logger: logr.Discard(),
	})
	assert.NoError(t, err)
}
