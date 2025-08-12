package redis

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
    rdb *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
    return &RateLimiter{rdb: rdb}
}

// CheckLimit uses fixed window with INCR + EXPIRE
func (rl *RateLimiter) CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    count, err := rl.rdb.Incr(ctx, key).Result()
    if err != nil {
        return false, err
    }
    if count == 1 {
        _ = rl.rdb.PExpire(ctx, key, window).Err()
    }
    return count <= int64(limit), nil
}