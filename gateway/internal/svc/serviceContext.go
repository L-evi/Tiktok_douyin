package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"train-tiktok/gateway/internal/config"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/identity/types/identity"
)

type ServiceContext struct {
	Config      config.Config
	IdentityRpc identity.IdentityClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		IdentityRpc: identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpc)),
	}
}
