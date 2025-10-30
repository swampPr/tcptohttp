package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParser(t *testing.T) {
	// TESTING: Valid single headers
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:     barbar\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("HOST"))
	assert.Equal(t, "barbar", headers.Get("FoOfOo"))
	assert.Equal(t, "", headers.Get("missingkey"))
	assert.Equal(t, 45, n)
	assert.True(t, done)

	// TESTING: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// TESTING: Invalid header name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}
