package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
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
	Config      config.Config
	IdentityRpc identity.IdentityClient
	VideoRpc    video.VideoClient
	UserRpc     user.UserClient
	Auth        rest.Middleware
	AuthPass    rest.Middleware
	PublicPath  string
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 视频临时存储路径
	_publicPath := os.Getenv("PUBLIC_BASE_PATH")
	if _publicPath == "" {
		_publicPath = c.PublicPath
	}

	// isdebug
	if isDebug, ok := os.LookupEnv("DEBUG"); ok {
		if isDebug == "true" {
			c.Log.Level = "debug"
			c.Log.Mode = "console"
		} else {
			c.Log.Level = "info"
			c.Log.Mode = "file"
			c.Log.KeepDays = 60
			c.Log.Rotation = "daily"
			c.Log.Encoding = "json"
		}
	}
	logx.MustSetup(c.Log)

	return &ServiceContext{
		Config:      c,
		IdentityRpc: identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpc)),
		VideoRpc:    videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Auth:        middleware.NewAuthMiddleware(c.IdentityRpc).Handle,
		AuthPass:    middleware.NewAuthPassMiddleware(c.IdentityRpc).Handle,
		PublicPath:  _publicPath,
	}
}
