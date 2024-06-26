package fx

import (
	"fmt"

	"github.com/xyzbit/gpkg/threading"
	"github.com/zeromicro/go-zero/core/logx"
)

// Pipeline takes a list of functions and returns a function whose param will be passed into
// the functions one by one.
func Pipeline[T any](funcs ...func(T) T) func(T) T {
	return func(arg T) (result T) {
		result = arg
		for _, fn := range funcs {
			result = fn(result)
		}
		return
	}
}

func PipelineVoid[T any](funcs ...func(T)) func(T) {
	return func(arg T) {
		for _, fn := range funcs {
			fn(arg)
		}
		return
	}
}

func PipelineE[T any](funcs ...func(T) error) func(T) {
	return func(arg T) {
		for _, fn := range funcs {
			if err := fn(arg); err != nil {
				fmt.Printf("pipeline exec failed: %+v", err)
				return
			}
		}
		return
	}
}

func PipelineAsyncVoid[T any](funcs ...func(T)) func(<-chan T) {
	funcEs := make([]func(T) error, 0, len(funcs))
	for _, fn := range funcs {
		funcEs = append(funcEs, func(t T) error {
			fn(t)
			return nil
		})
	}
	return PipelineAsyncE(funcEs...)
}

// Pipeline takes a list of functions and returns a function whose param will be passed into
// async executed between each one.
func PipelineAsyncE[T any](funcs ...func(T) error) func(<-chan T) {
	if len(funcs) == 0 {
		return func(<-chan T) {}
	}
	if len(funcs) == 1 {
		return func(c <-chan T) {
			for args := range c {
				if err := funcs[0](args); err != nil {
					logx.Error(err)
					return
				}
			}
		}
	}

	return func(argsChan <-chan T) {
		chans := make([]chan T, len(funcs)-1)
		for i := 0; i < len(chans); i++ {
			chans[i] = make(chan T)
		}
		wait := make(chan struct{})
		errChan := make(chan error)

		for i := 1; i < len(funcs); i++ {
			i := i
			threading.GoSafe(func() {
				lastFunc := (i == len(funcs)-1)

				for args := range chans[i-1] {
					if err := funcs[i](args); err != nil {
						errChan <- err
					}
					if !lastFunc {
						chans[i] <- args
					}
				}

				if lastFunc {
					wait <- struct{}{}
				} else {
					close(chans[i])
				}
			})
		}

		threading.GoSafe(func() {
			for args := range argsChan {
				if err := funcs[0](args); err != nil {
					errChan <- err
					return
				}
				chans[0] <- args
			}
			close(chans[0])
		})

		for {
			select {
			case e := <-errChan:
				fmt.Printf("pipeline exec error: %+v", e)
				return
			case <-wait:
				return
			}
		}
	}
}
