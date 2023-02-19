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
	IdentityRpcConf zrpc.RpcClientConf
	RedisConf       redisutil.RedisConf
}
