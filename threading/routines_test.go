package threading

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
)

func TestRunSafe(t *testing.T) {
	log.SetOutput(io.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan lang.PlaceholderType)
	go RunSafe(func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRunSafeE(t *testing.T) {
	log.SetOutput(io.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan lang.PlaceholderType)
	_ = RunSafeE(func() error {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}
