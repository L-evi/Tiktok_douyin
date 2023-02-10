package redisutil

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// type RedisConf
type RedisConf struct {
	Addr        string
	Password    string
	DB          int
	MaxIdle     int
	MinIdle     int
	PoolSize    int
	MaxLifeTime int
}

// New xx

func New(conf RedisConf) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.DB,
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
