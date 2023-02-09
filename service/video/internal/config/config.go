package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	StorageBaseUrl StorageStruct
}

type StorageStruct struct {
	Local string
	Cos   string
}
