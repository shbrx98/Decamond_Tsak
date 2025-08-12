package redis

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "time"

    "github.com/redis/go-redis/v9"
)

type OTPStore struct {
    rdb    *redis.Client
    pepper []byte
}

func NewOTPStore(rdb *redis.Client, pepper string) *OTPStore {
    return &OTPStore{rdb: rdb, pepper: []byte(pepper)}
}

func (s *OTPStore) key(phone string) string      { return "otp:" + phone }
func (s *OTPStore) attemptsKey(p string) string  { return "otp_attempts:" + p }

func (s *OTPStore) Save(ctx context.Context, phone, code string, ttl time.Duration) error {
    hashed := s.hash(phone, code)
    if err := s.rdb.Set(ctx, s.key(phone), hashed, ttl).Err(); err != nil {
        return err
    }
    // reset attempts alongside OTP creation
    _ = s.rdb.Del(ctx, s.attemptsKey(phone)).Err()
    return nil
}

func (s *OTPStore) GetHash(ctx context.Context, phone string) (string, error) {
    return s.rdb.Get(ctx, s.key(phone)).Result()
}

func (s *OTPStore) Delete(ctx context.Context, phone string) error {
    return s.rdb.Del(ctx, s.key(phone)).Err()
}

func (s *OTPStore) IncAttempts(ctx context.Context, phone string, ttl time.Duration) (int64, error) {
    k := s.attemptsKey(phone)
    c, err := s.rdb.Incr(ctx, k).Result()
    if err != nil {
        return 0, err
    }
    if c == 1 {
        _ = s.rdb.PExpire(ctx, k, ttl).Err()
    }
    return c, nil
}

func (s *OTPStore) ResetAttempts(ctx context.Context, phone string) error {
    return s.rdb.Del(ctx, s.attemptsKey(phone)).Err()
}

func (s *OTPStore) hash(phone, code string) string {
    mac := hmac.New(sha256.New, s.pepper)
    mac.Write([]byte(phone))
    mac.Write([]byte("|"))
    mac.Write([]byte(code))
    return hex.EncodeToString(mac.Sum(nil))
}