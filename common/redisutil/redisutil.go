package redisutil

import (
	"context"
	"github.com/redis/go-redis/v9"

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
	//  check if connect
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panicf("redis init failed, err:%v", err)
	}
	return rdb
}
