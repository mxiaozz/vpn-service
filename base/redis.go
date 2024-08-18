package base

import (
	"context"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"vpn-web.funcworks.net/gb"
)

func initRedis() {
	// redisCfg := gb.Config.Redis
	// var client redis.UniversalClient
	// // 使用集群模式
	// if redisCfg.UseCluster {
	// 	client = redis.NewClusterClient(&redis.ClusterOptions{
	// 		Addrs:        redisCfg.ClusterAddrs,
	// 		Password:     redisCfg.Password,
	// 		DialTimeout:  5 * time.Second,
	// 		ReadTimeout:  30 * time.Second,
	// 		WriteTimeout: 30 * time.Second,
	// 	})
	// } else {
	// 	// 使用单例模式
	// 	client = redis.NewClient(&redis.Options{
	// 		Addr:         redisCfg.Addr,
	// 		Password:     redisCfg.Password,
	// 		DB:           redisCfg.DB,
	// 		DialTimeout:  5 * time.Second,
	// 		ReadTimeout:  30 * time.Second,
	// 		WriteTimeout: 30 * time.Second,
	// 	})
	// }
	// gb.RedisClient = client
	gb.RedisClient = newLocalCache()

	gb.Logger.Debug("initialized redis client")
}

type localCache struct {
	redis.UniversalClient
	cache *cache.Cache
}

func newLocalCache() *localCache {
	return &localCache{cache: cache.New(0, 10*time.Minute)}
}

func (lc *localCache) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if value, exist := lc.cache.Get(key); exist {
		cmd.SetVal(value.(string))
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

func (lc *localCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	lc.cache.Set(key, value, expiration)
	return redis.NewStatusCmd(ctx)
}

func (lc *localCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	for _, key := range keys {
		lc.cache.Delete(key)
	}
	return redis.NewIntCmd(ctx)
}

func (lc *localCache) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	cmd := redis.NewStringSliceCmd(ctx)
	if pattern == "" {
		cmd.SetVal([]string{})
		return cmd
	}
	pattern = strings.TrimSuffix(pattern, "*")

	data := lc.cache.Items()
	keys := make([]string, 0, len(data))
	for key := range data {
		if pattern == "" || strings.HasPrefix(key, pattern) {
			keys = append(keys, key)
		}
	}
	cmd.SetVal(keys)

	return cmd
}
