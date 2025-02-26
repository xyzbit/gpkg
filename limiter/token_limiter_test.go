package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestTokenLimiter_Allow(t *testing.T) {
	s := miniredis.RunT(t)
	cli := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	tl, err := NewTokenLimiter(cli, "", 120)
	if err != nil {
		t.Error(err)
	}

	var (
		wg       sync.WaitGroup
		start    = time.Now()
		reqCount atomic.Int32
	)

	loopCount := 200
	allowCount := 5
	wg.Add(loopCount)
	go func() {
		for i := 0; i < loopCount; i++ {
			go func(idx int) {
				defer wg.Done()
				if allow, _ := tl.AllowN(allowCount); allow {
					reqCount.Add(int32(allowCount))
					return
				}
			}(i)
		}
	}()

	loopCount = 200
	allowCount = 1
	wg.Add(loopCount)
	go func() {
		for i := 1000; i < loopCount+1000; i++ {
			go func(idx int) {
				defer wg.Done()
				if allow, _ := tl.AllowN(allowCount); allow {
					reqCount.Add(int32(allowCount))
					return
				}
			}(i)
		}
	}()

	wg.Wait()

	t.Logf("since: %v, reqcount: %v", time.Since(start), reqCount.Load())
	assert.Equal(t, tl.QPS, int(reqCount.Load()))

	// gotEndTime := time.Now()
	// wantEndTime := start.Add(time.Duration(reqCount.Load()/int32(tl.QPS)) * time.Second)
	// mistake := 10 * time.Millisecond // 误差容忍
	// assert.WithinRange(t, gotEndTime, wantEndTime.Add(-10*mistake), wantEndTime.Add(10*mistake))
}

func TestTokenLimiter_DelayN(t *testing.T) {
	s := miniredis.RunT(t)
	cli := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	tl, err := NewTokenLimiter(cli, "", 10)
	if err != nil {
		t.Error(err)
	}

	var (
		tryTimes = 10000
		start    = time.Now()
		wg       sync.WaitGroup
		reqCount atomic.Int32
	)

	for i := 0; i < tryTimes; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				delay, err := tl.DelayN(1)
				if err != nil {
					t.Error(err)
				}
				if delay == 0 {
					reqCount.Add(1)
					return
				}
				time.Sleep(delay)
			}
		}(i)
	}

	wg.Wait()

	gotEndTime := time.Now()
	wantEndTime := start.Add(time.Duration(reqCount.Load()/int32(tl.QPS)) * time.Second)
	mistake := 1 * time.Second // 误差容忍, 随着请求时间的增长, 误差会变大
	t.Logf("real mistake:%v, req count:%d", wantEndTime.Sub(gotEndTime), reqCount.Load())
	assert.WithinRange(t, gotEndTime, wantEndTime.Add(-1*mistake), wantEndTime.Add(1*mistake))
}

func BenchmarkTokenLimiter_AllowN(b *testing.B) {
	s := miniredis.RunT(b)
	cli := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	tl, err := NewTokenLimiter(cli, "", 100)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		go func(idx int) {
			if allow, _ := tl.AllowN(10); allow {
				// fmt.Printf("No.%d execute success  ---at time:[%s] \n", idx, time.Now().Format("2006-01-02 15:04:05"))
				return
			}
		}(i)
	}
}
