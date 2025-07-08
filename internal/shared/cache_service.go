package shared

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd

	Incr(ctx context.Context, key string) *redis.IntCmd
	IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd
	Decr(ctx context.Context, key string) *redis.IntCmd
	DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd

	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LLen(ctx context.Context, key string) *redis.IntCmd

	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd

	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub

	ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd
	ZRemRangeByScore(ctx context.Context, key string, min, max string) *redis.IntCmd
	ZCard(ctx context.Context, key string) *redis.IntCmd
	ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd

	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd

	TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
	Close() error
}

func NewRedisCacheService(addr string, password string, tlsEnabled bool) CacheService {
	if tlsEnabled {
		return redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}
