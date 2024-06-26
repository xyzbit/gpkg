package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperFirst(t *testing.T) {
	cases := map[string]string{
		"":     "",
		"foo":  "Foo",
		"bAR":  "BAR",
		"FOo":  "FOo",
		"fOo大": "FOo大",
	}

	for k, v := range cases {
		assert.Equal(t, v, UpperFirst(k))
	}
}

func TestLowerFirst(t *testing.T) {
	cases := map[string]string{
		"":     "",
		"foo":  "foo",
		"bAR":  "bAR",
		"FOo":  "fOo",
		"fOo大": "fOo大",
	}

	for k, v := range cases {
		assert.Equal(t, v, LowerFirst(k))
	}
}
