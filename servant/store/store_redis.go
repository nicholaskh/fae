package store

import (
	"github.com/nicholaskh/fae/config"
	"github.com/nicholaskh/fae/servant/redis"
)

type RedisStore struct {
	pool  string
	redis *redis.Client
}

func NewRedisStore(pool string, cf *config.ConfigRedis) *RedisStore {
	this := &RedisStore{pool: pool, redis: redis.New(cf)}
	return this
}

func (this *RedisStore) Get(key string) (val interface{}, present bool) {
	var err error
	if val, err = this.redis.Get(this.pool, key); err == nil {
		present = true
	}
	return
}

func (this *RedisStore) Set(key string, val interface{}) {
	this.redis.Set(this.pool, key, val)
}

func (this *RedisStore) Del(key string) {
	this.redis.Del(this.pool, key)
}
