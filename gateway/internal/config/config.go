package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	IdentityRpc zrpc.RpcClientConf
	VideoRpc    zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	PublicPath  string
	ChatRpc     zrpc.RpcClientConf
	Cos         CosConf
}

type CosConf struct {
	Enable    bool
	SecretId  string
	SecretKey string
	BucketUrl string
	Path      string
}
