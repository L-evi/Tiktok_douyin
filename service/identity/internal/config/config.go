package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Jwt struct {
		SigningKey string
	}
	Conf struct {
		GravatarBaseURL   string
		DefaultBackground string
	}
}
