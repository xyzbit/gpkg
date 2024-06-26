package fx

import (
	"fmt"
	"testing"
	"time"
)

func ExamplePipelineAsync() {
	pipline := PipelineAsyncE(
		func(t *testPipline) error {
			t.name = "testname"
			fmt.Println(t.name)
			return nil
		}, func(t *testPipline) error {
			t.args = "testargs"
			fmt.Println(t.args)
			return nil
		},
	)

	pipline(sliceToChannel(3, []*testPipline{{}}))

	// Output:
	// testname
	// testargs
}

func ExamplePipelineAsyncUseTime() {
	start := time.Now()
	pipline := PipelineAsyncE(
		func(t *testPipline) error {
			t.name = "testname"
			time.Sleep(time.Second * 1)
			return nil
		}, func(t *testPipline) error {
			t.args = "testargs"
			time.Sleep(time.Second * 1)
			return nil
		},
	)

	pipline(sliceToChannel(3, []*testPipline{{}, {}, {}}))
	fmt.Println(int(time.Since(start).Seconds()))

	// Output:
	// 4
}

type testPipline struct {
	name string
	args string
}

// go test -timeout 30s -run ^TestPipelineAsyncE$ devops/pkg/utils/fuction -v -count=1 -race
func TestPipelineAsyncE(t *testing.T) {
	tests := []struct {
		name  string
		funcs []func(*testPipline) error
	}{
		{
			name: "test1",
			funcs: []func(*testPipline) error{
				func(t *testPipline) error {
					t.name = "testname"
					return nil
				},
				func(t *testPipline) error {
					t.args = "testargs"
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipline := PipelineAsyncE(tt.funcs...)
			pipline(sliceToChannel(3, []*testPipline{
				{}, {}, {},
			}))
		})
	}
}

// go test -timeout 30s -run ^TestPipelineAsyncEWithError$ devops/pkg/utils/fuction -v -count=1 -race
func TestPipelineAsyncEWithError(t *testing.T) {
	tests := []struct {
		name  string
		funcs []func(*testPipline) error
	}{
		{
			name: "test1",
			funcs: []func(*testPipline) error{
				func(t *testPipline) error {
					if t.name == "hasherr" {
						return fmt.Errorf("testname error")
					}
					return nil
				},
				func(t *testPipline) error {
					t.args = "testargs"
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipline := PipelineAsyncE(tt.funcs...)
			pipline(sliceToChannel(3, []*testPipline{
				{}, {name: "hasherr"}, {},
			}))
		})
	}
}

func sliceToChannel[T any](bufferSize int, collection []T) <-chan T {
	ch := make(chan T, bufferSize)

	go func() {
		for _, item := range collection {
			ch <- item
		}

		close(ch)
	}()

	return ch
}
