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
	Comment               = video.Comment
	CommentActionReq      = video.CommentActionReq
	CommentActionResp     = video.CommentActionResp
	CommentCountReq       = video.CommentCountReq
	CommentCountResp      = video.CommentCountResp
	CommentListReq        = video.CommentListReq
	CommentListResp       = video.CommentListResp
	FavoriteActionReq     = video.FavoriteActionReq
	FavoriteActionResp    = video.FavoriteActionResp
	FavoriteCountReq      = video.FavoriteCountReq
	FavoriteCountResp     = video.FavoriteCountResp
	FavoriteListReq       = video.FavoriteListReq
	FavoriteListResp      = video.FavoriteListResp
	FavoritedCountReq     = video.FavoritedCountReq
	FavoritedCountResp    = video.FavoritedCountResp
	FeedReq               = video.FeedReq
	FeedResp              = video.FeedResp
	GetVideoByHashReq     = video.GetVideoByHashReq
	GetVideoByHashResp    = video.GetVideoByHashResp
	IsFavoriteReq         = video.IsFavoriteReq
	IsFavoriteResp        = video.IsFavoriteResp
	PublishListReq        = video.PublishListReq
	PublishListResp       = video.PublishListResp
	PublishReq            = video.PublishReq
	PublishResp           = video.PublishResp
	Resp                  = video.Resp
	UserFavoriteCountReq  = video.UserFavoriteCountReq
	UserFavoriteCountResp = video.UserFavoriteCountResp
	VideoX                = video.VideoX
	WorkCountReq          = video.WorkCountReq
	WorkCountResp         = video.WorkCountResp

	Video interface {
		Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error)
		PublishList(ctx context.Context, in *PublishListReq, opts ...grpc.CallOption) (*PublishListResp, error)
		Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error)
		CommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error)
		CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error)
		FavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error)
		FavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error)
		FavoriteCount(ctx context.Context, in *FavoriteCountReq, opts ...grpc.CallOption) (*FavoriteCountResp, error)
		CommentCount(ctx context.Context, in *CommentCountReq, opts ...grpc.CallOption) (*CommentCountResp, error)
		IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error)
		WorkCount(ctx context.Context, in *WorkCountReq, opts ...grpc.CallOption) (*WorkCountResp, error)
		FavoritedCount(ctx context.Context, in *FavoritedCountReq, opts ...grpc.CallOption) (*FavoritedCountResp, error)
		UserFavoriteCount(ctx context.Context, in *UserFavoriteCountReq, opts ...grpc.CallOption) (*UserFavoriteCountResp, error)
		GetVideoByHash(ctx context.Context, in *GetVideoByHashReq, opts ...grpc.CallOption) (*GetVideoByHashResp, error)
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

func (m *defaultVideo) PublishList(ctx context.Context, in *PublishListReq, opts ...grpc.CallOption) (*PublishListResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.PublishList(ctx, in, opts...)
}

func (m *defaultVideo) Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.Feed(ctx, in, opts...)
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

func (m *defaultVideo) FavoriteCount(ctx context.Context, in *FavoriteCountReq, opts ...grpc.CallOption) (*FavoriteCountResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.FavoriteCount(ctx, in, opts...)
}

func (m *defaultVideo) CommentCount(ctx context.Context, in *CommentCountReq, opts ...grpc.CallOption) (*CommentCountResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.CommentCount(ctx, in, opts...)
}

func (m *defaultVideo) IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.IsFavorite(ctx, in, opts...)
}

func (m *defaultVideo) WorkCount(ctx context.Context, in *WorkCountReq, opts ...grpc.CallOption) (*WorkCountResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.WorkCount(ctx, in, opts...)
}

func (m *defaultVideo) FavoritedCount(ctx context.Context, in *FavoritedCountReq, opts ...grpc.CallOption) (*FavoritedCountResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.FavoritedCount(ctx, in, opts...)
}

func (m *defaultVideo) UserFavoriteCount(ctx context.Context, in *UserFavoriteCountReq, opts ...grpc.CallOption) (*UserFavoriteCountResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.UserFavoriteCount(ctx, in, opts...)
}

func (m *defaultVideo) GetVideoByHash(ctx context.Context, in *GetVideoByHashReq, opts ...grpc.CallOption) (*GetVideoByHashResp, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.GetVideoByHash(ctx, in, opts...)
}
