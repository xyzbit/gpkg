package threading

import (
	"fmt"
	"runtime/debug"
)

func RunSafe(fn func()) {
	defer func() {
		if result := recover(); result != nil {
			stackInfo := string(debug.Stack())
			fmt.Println("RunSafe_panic: ", stackInfo)
		}
	}()

	fn()
}

func RunSafeE(fn func() error) error {
	defer func() {
		if result := recover(); result != nil {
			stackInfo := string(debug.Stack())
			fmt.Println("RunSafeE_panic: ", stackInfo)
		}
	}()

	return fn()
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}
