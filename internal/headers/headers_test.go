package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:      barbar      \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	hostValue, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", hostValue)
	fooValue, _ := headers.Get("FooFoo")
	assert.Equal(t, "barbar", fooValue)
	missingValue, ok := headers.Get("MissingKey")
	assert.Equal(t, "", missingValue)
	assert.False(t, ok)
	assert.Equal(t, 52, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n)")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42069\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	combinedHostValue, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069, localhost:42069", combinedHostValue)
	assert.False(t, done)
}
