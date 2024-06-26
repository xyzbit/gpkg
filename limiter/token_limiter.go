package limiter

import (
	"context"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	defaultLimiterKey = "default"
)

var (
	ErrQPSNonPositive = errors.New("qps is non-positive")
	ErrExecLua        = errors.New("exec lua")
)

type TokenLimiter struct {
	luaSha string
	key    string

	redisClient *redis.Client
	QPS         int
}

func NewTokenLimiter(cli *redis.Client, key string, qps int) (*TokenLimiter, error) {
	ctx := context.Background()
	err := cli.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	if key == "" {
		key = defaultLimiterKey
	}
	key += "_token_limiter"

	return &TokenLimiter{
		luaSha:      cli.ScriptLoad(ctx, tokenLimiterLuaScript).Val(),
		key:         key,
		redisClient: cli,
		QPS:         qps,
	}, nil
}

func (t *TokenLimiter) Allow() (bool, error) {
	return t.AllowN(1)
}

func (t *TokenLimiter) AllowN(n int) (bool, error) {
	if t.QPS <= 0 {
		return false, ErrQPSNonPositive
	}

	maxToken := t.QPS
	needToken := n
	tokenPerSec := t.QPS

	raw, err := t.redisClient.EvalSha(context.Background(), t.luaSha, []string{t.key}, needToken, maxToken, tokenPerSec, time.Now().Unix()).Result()
	if err != nil {
		return false, errors.Wrap(ErrExecLua, err.Error())
	}

	tokens, ok := raw.(int64)
	if !ok || tokens == 0 || tokens != int64(n) {
		return false, nil
	}

	return true, nil
}

func (t *TokenLimiter) DelayN(n int) (time.Duration, error) {
	allow, err := t.AllowN(n)
	if err != nil {
		return time.Duration(math.MaxInt64), err
	}
	if allow {
		return 0, nil
	}

	return time.Duration(float64(n)/float64(t.QPS)*1e9) * time.Nanosecond, nil
}
