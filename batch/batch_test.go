package batch

import (
	"fmt"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	stopCh := make(chan struct{})

	bt := NewBatch(func(ts []int) {
		fmt.Println(len(ts))
	},
		WithBatchInterval[int](2000),
		WithBatchSize[int](10),
		WithDone[int](stopCh))

	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			bt.SendData(1)
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		close(stopCh)
	}()

	select {
	case <-stopCh:
		fmt.Println("done")
	}
}
