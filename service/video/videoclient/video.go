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
	Comment            = video.Comment
	CommentActionReq   = video.CommentActionReq
	CommentActionResp  = video.CommentActionResp
	CommentListReq     = video.CommentListReq
	CommentListResp    = video.CommentListResp
	FavoriteActionReq  = video.FavoriteActionReq
	FavoriteActionResp = video.FavoriteActionResp
	FavoriteListReq    = video.FavoriteListReq
	FavoriteListResp   = video.FavoriteListResp
	PublishReq         = video.PublishReq
	PublishResp        = video.PublishResp
	Resp               = video.Resp
	User               = video.User
	Video              = video.Video

	Video interface {
		Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error)
		CommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error)
		CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error)
		FavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error)
		FavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error)
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

func (m *defaultVideo) CommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.CommentAction(ctx, in, opts...)
}

func (m *defaultVideo) CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.CommentList(ctx, in, opts...)
}

func (m *defaultVideo) FavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.FavoriteAction(ctx, in, opts...)
}

func (m *defaultVideo) FavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.FavoriteList(ctx, in, opts...)
}
