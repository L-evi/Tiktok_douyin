package svc

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
	"net/http"
	"net/url"
	"os"
	"train-tiktok/common/logset"
	"train-tiktok/gateway/internal/config"
	"train-tiktok/gateway/internal/middleware"
	"train-tiktok/service/chat/chatclient"
	"train-tiktok/service/chat/types/chat"
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
	ChatRpc     chat.ChatClient
}

func NewServiceContext(c config.Config) *ServiceContext {

	var debug = false
	if isDebug, ok := os.LookupEnv("DEBUG"); ok {
		if isDebug == "true" {
			debug = true
		}
	}
	logset.Handler(debug, c.Log)

	// 视频临时存储路径
	if publicPath, ok := os.LookupEnv("PUBLIC_BASE_PATH"); ok {
		c.PublicPath = publicPath
	}

	if etcdEndpoint, ok := os.LookupEnv("ETCD_ENDPOINT"); ok {
		c.IdentityRpc.Etcd.Hosts = []string{etcdEndpoint}
		c.VideoRpc.Etcd.Hosts = []string{etcdEndpoint}
		c.UserRpc.Etcd.Hosts = []string{etcdEndpoint}
		c.ChatRpc.Etcd.Hosts = []string{etcdEndpoint}
	}

	// cos
	if CosBucketEnabled, ok := os.LookupEnv("COS_BUCKET_ENABLE"); ok {
		c.Cos.Enable = false
		if CosBucketEnabled == "true" {
			c.Cos.Enable = true
		}
	}
	if CosBucketUrl, ok := os.LookupEnv("COS_BUCKET_URL"); ok {
		c.Cos.BucketUrl = CosBucketUrl
	}
	if CosSecretId, ok := os.LookupEnv("COS_SECRET_ID"); ok {
		c.Cos.SecretId = CosSecretId
	}
	if CosSecretKey, ok := os.LookupEnv("COS_SECRET_KEY"); ok {
		c.Cos.SecretKey = CosSecretKey
	}
	if CosPath, ok := os.LookupEnv("COS_PATH"); ok {
		c.Cos.Path = CosPath
	}

	// check Bucket exist
	if c.Cos.Enable {
		bucketURL, _ := url.Parse(c.Cos.BucketUrl)
		b := &cos.BaseURL{BucketURL: bucketURL}

		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  c.Cos.SecretId,
				SecretKey: c.Cos.SecretKey,
			},
		})

		ok, err := client.Bucket.IsExist(context.Background())

		if err != nil {
			log.Panicf("Bucket exists Check Failed, %s", err)
		} else if !ok {
			log.Panicf("Bucket not exists")
		}
	}

	return &ServiceContext{
		Config:      c,
		IdentityRpc: identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpc)),
		VideoRpc:    videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ChatRpc:     chatclient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		Auth:        middleware.NewAuthMiddleware(c.IdentityRpc).Handle,
		AuthPass:    middleware.NewAuthPassMiddleware(c.IdentityRpc).Handle,
	}
}
