// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package videoclient

import (
	"context"

	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	PublishReq  = video.PublishReq
	PublishResp = video.PublishResp
	Resp        = video.Resp

	Video interface {
		Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error)
	}

	defaultVideo struct {
		cli zrpc.Client
	}
)

func NewVideo(cli zrpc.Client) Video {
	return &defaultVideo{
		cli: cli,
	}
}

func (m *defaultVideo) Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.Publish(ctx, in, opts...)
}