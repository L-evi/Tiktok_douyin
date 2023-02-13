// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package videos

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
	CommentCountReq    = video.CommentCountReq
	CommentCountResp   = video.CommentCountResp
	CommentListReq     = video.CommentListReq
	CommentListResp    = video.CommentListResp
	FavoriteActionReq  = video.FavoriteActionReq
	FavoriteActionResp = video.FavoriteActionResp
	FavoriteCountReq   = video.FavoriteCountReq
	FavoriteCountResp  = video.FavoriteCountResp
	FavoriteListReq    = video.FavoriteListReq
	FavoriteListResp   = video.FavoriteListResp
	FavoriteVideo      = video.FavoriteVideo
	FeedReq            = video.FeedReq
	FeedResp           = video.FeedResp
	IsFavoriteReq      = video.IsFavoriteReq
	IsFavoriteResp     = video.IsFavoriteResp
	PublishReq         = video.PublishReq
	PublishResp        = video.PublishResp
	Resp               = video.Resp
	Video              = video.Video

	Videos interface {
		Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error)
		Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error)
		CommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error)
		CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error)
		FavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error)
		FavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error)
		FavoriteCount(ctx context.Context, in *FavoriteCountReq, opts ...grpc.CallOption) (*FavoriteCountResp, error)
		CommentCount(ctx context.Context, in *CommentCountReq, opts ...grpc.CallOption) (*CommentCountResp, error)
		IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error)
	}

	defaultVideos struct {
		cli zrpc.Client
	}
)

func NewVideos(cli zrpc.Client) Videos {
	return &defaultVideos{
		cli: cli,
	}
}

func (m *defaultVideos) Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.Publish(ctx, in, opts...)
}

func (m *defaultVideos) Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.Feed(ctx, in, opts...)
}

func (m *defaultVideos) CommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.CommentAction(ctx, in, opts...)
}

func (m *defaultVideos) CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.CommentList(ctx, in, opts...)
}

func (m *defaultVideos) FavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.FavoriteAction(ctx, in, opts...)
}

func (m *defaultVideos) FavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.FavoriteList(ctx, in, opts...)
}

func (m *defaultVideos) FavoriteCount(ctx context.Context, in *FavoriteCountReq, opts ...grpc.CallOption) (*FavoriteCountResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.FavoriteCount(ctx, in, opts...)
}

func (m *defaultVideos) CommentCount(ctx context.Context, in *CommentCountReq, opts ...grpc.CallOption) (*CommentCountResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.CommentCount(ctx, in, opts...)
}

func (m *defaultVideos) IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error) {
	client := video.NewVideosClient(m.cli.Conn())
	return client.IsFavorite(ctx, in, opts...)
}