package gb

import (
	"context"
	"time"
)

type redisProxy struct {
}

func (rc *redisProxy) Get(key string) (string, error) {
	return RedisClient.Get(context.Background(), key).Result()
}

func (rc *redisProxy) GetInt(key string) (int, error) {
	return RedisClient.Get(context.Background(), key).Int()
}

func (rc *redisProxy) Set(key string, value any) error {
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

// expiration 单位：秒
func (rc *redisProxy) SetEx(key string, value any, expiration int) error {
	return RedisClient.Set(context.Background(), key, value, time.Duration(expiration)*time.Second).Err()
}

func (rc *redisProxy) Delete(key string) error {
	return RedisClient.Del(context.Background(), key).Err()
}
