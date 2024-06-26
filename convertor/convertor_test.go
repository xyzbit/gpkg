package convertor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes2Str(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{
			name:  "Test Case 1",
			input: "Hello, World!",
			want:  []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33},
		},
		{
			name:  "Test Case 2",
			input: "Testing",
			want:  []byte{84, 101, 115, 116, 105, 110, 103},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Str2Bytes(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Str2Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStr2Bytes(t *testing.T) {
	b := []byte("Hello, World!")
	expected := "Hello, World!"

	result := Bytes2Str(b)

	if result != expected {
		t.Errorf("Got %q; want %q", result, expected)
	}
}

func TestToString(t *testing.T) {
	var x interface{}
	x = 123 // 赋值为整数
	assert.Equal(t, "123", ToString(x))
	x = "this is test"
	assert.Equal(t, "this is test", ToString(x))
	x = []byte("this is test")
	assert.Equal(t, "this is test", ToString(x))

	x = User{Username: "lala", Age: 25}
	assert.Equal(t, "{\"Username\":\"lala\",\"Age\":25,\"Address\":{\"Detail\":\"\"}}", ToString(x))

	x = map[string]string{
		"name": "lala",
	}
	assert.Equal(t, "{\"name\":\"lala\"}", ToString(x))
}

type User struct {
	Username string
	Age      int
	Address  Address
}
type Address struct {
	Detail string
}
