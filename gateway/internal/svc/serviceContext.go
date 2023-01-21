package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"train-tiktok/gateway/internal/config"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/identity/types/identity"
)

type ServiceContext struct {
	Config       config.Config
	IdentityRpc  identity.IdentityClient
	VideoTmpPath string
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 视频临时存储路径
	_videoTmpPath := os.Getenv("VIDEO_TMP_PATH")
	if _videoTmpPath == "" {
		_videoTmpPath = c.VideoTmpPath
	}

	return &ServiceContext{
		Config:       c,
		IdentityRpc:  identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpc)),
		VideoTmpPath: _videoTmpPath,
	}
}
