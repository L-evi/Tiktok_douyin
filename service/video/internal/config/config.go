package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr        string
		Password    string
		DB          int
		MinIdle     int
		PoolSize    int
		MaxLifeTime int
	}
	IdentityRpcConf zrpc.RpcClientConf
	StorageBaseUrl  StorageStruct
}

type StorageStruct struct {
	Local string
	Cos   string
}
