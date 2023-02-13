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
	RedisConf       redisutil.RedisConf
	IdentityRpcConf zrpc.RpcClientConf
	StorageBaseUrl  StorageStruct
}

type StorageStruct struct {
	Local string
	Cos   string
}
