package redisutil

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// RedisConf struct
type RedisConf struct {
	Addr        string
	Passwd      string
	Db          int
	MinIdle     int
	PoolSize    int
	MaxLifeTime int
	Prefix      string
}

// New xx

func New(conf RedisConf) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Passwd,
		DB:           conf.Db,
		MinIdleConns: conf.MinIdle,
		MaxRetries:   3,
		PoolSize:     conf.PoolSize,
	})
	var ctx context.Context
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Panicf("redis init failed, err:%v", err)
	}
	return rdb
}
