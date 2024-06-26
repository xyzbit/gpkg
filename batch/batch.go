package batch

import (
	"fmt"
	"time"

	"github.com/xyzbit/pkg/threading"
)

const (
	defaultBatchSize     = 500  // 默认每批数据大小：500
	defaultBatchInterval = 2000 // 默认上传时间间隔：2000 毫秒
)

type Option[T any] func(b *Batch[T])

type Batch[T any] struct {
	dataCh         chan T        // 数据上传队列
	batchSize      int           // 打包上传的数量，默认 100，累积达到 batchSize 时立即上报，否则每 batchSeconds 秒上报一次
	batchInterval  int           // 上传间隔时间，单位毫秒
	tryTimes       int           // 上传时候尝试次数
	tryInterval    int           // 每次尝试之间的时间间隔，单位毫秒
	batchProcessor func([]T)     // 批处理函数
	done           chan struct{} // 停止控制
}

// WithBatchSize 批量数量
func WithBatchSize[T any](size int) Option[T] {
	return func(b *Batch[T]) {
		if size > 0 {
			b.batchSize = size
		}
	}
}

// WithBatchInterval 处理的间隔时间，单位毫秒
func WithBatchInterval[T any](interval int) Option[T] {
	return func(b *Batch[T]) {
		if interval > 0 {
			b.batchInterval = interval
		}
	}
}

// WithDone 添加停止控制
func WithDone[T any](done chan struct{}) Option[T] {
	return func(b *Batch[T]) {
		b.done = done
	}
}

// NewBatch 新建批处理对象
func NewBatch[T any](processor func([]T), opts ...Option[T]) *Batch[T] {
	b := &Batch[T]{
		dataCh:         make(chan T, 2000),
		batchSize:      defaultBatchSize,
		batchInterval:  defaultBatchInterval,
		batchProcessor: processor,
		done:           make(chan struct{}),
	}

	for _, opt := range opts {
		opt(b)
	}

	threading.GoSafe(b.start)
	return b
}

// SendData 发送数据到处理队列
func (b *Batch[T]) SendData(data T) {
	b.dataCh <- data
}

func (b *Batch[T]) start() {
	collector := make([]T, 0, b.batchSize)
	batchDuration := time.Duration(b.batchInterval) * time.Millisecond
	ticker := time.NewTicker(batchDuration)
	defer ticker.Stop()

	for {
		select {
		case <-b.done:
			fmt.Println("batch done")
			return
		case <-ticker.C:
			copiedCollector := make([]T, len(collector))
			copy(copiedCollector, collector)
			b.batchProcessor(copiedCollector)
			collector = make([]T, 0, b.batchSize)
		case msg := <-b.dataCh:
			collector = append(collector, msg)
			if len(collector) >= b.batchSize {
				copiedCollector := make([]T, len(collector))
				copy(copiedCollector, collector)
				b.batchProcessor(copiedCollector)
				collector = make([]T, 0, b.batchSize)
			}
		}
	}
}
