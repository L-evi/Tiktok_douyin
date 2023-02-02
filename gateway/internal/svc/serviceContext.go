package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"train-tiktok/gateway/internal/config"
	"train-tiktok/gateway/internal/middleware"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/identity/types/identity"
	"train-tiktok/service/user/types/user"
	"train-tiktok/service/user/userclient"
	"train-tiktok/service/video/types/video"
	"train-tiktok/service/video/videoclient"
)

type ServiceContext struct {
	Config       config.Config
	IdentityRpc  identity.IdentityClient
	VideoRpc     video.VideoClient
	UserRpc      user.UserClient
	VideoTmpPath string
	Auth         rest.Middleware
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
		VideoRpc:     videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:      userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoTmpPath: _videoTmpPath,
		Auth:         middleware.NewAuthMiddleware(c.IdentityRpc).Handle,
	}
}
