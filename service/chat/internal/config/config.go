package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"train-tiktok/common/redisutil"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Redis          redisutil.RedisConf
	StorageBaseUrl StorageStruct
}

type StorageStruct struct {
	Local string
	Cos   string
}
