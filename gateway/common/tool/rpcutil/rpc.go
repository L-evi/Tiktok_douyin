package rpcutil

import (
	"context"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	"train-tiktok/service/user/types/user"
	"train-tiktok/service/video/types/video"
)

func IsFavorite(c *svc.ServiceContext, ctx context.Context, userId int64, videoId int64) (bool, error) {
	var err error
	var resp *video.IsFavoriteResp
	if resp, err = c.VideoRpc.IsFavorite(ctx, &video.IsFavoriteReq{
		UserId:  userId,
		VideoId: videoId,
	}); err != nil {
		return false, err
	}
	return resp.IsFavorite, nil
}

func GetFavoriteCount(c *svc.ServiceContext, ctx context.Context, videoId int64) (int64, error) {
	var err error
	var resp *video.FavoriteCountResp
	if resp, err = c.VideoRpc.FavoriteCount(ctx, &video.FavoriteCountReq{
		VideoId: videoId,
	}); err != nil {
		return 0, err
	}
	return resp.FavoriteCount, nil
}

func GetCommentCount(c *svc.ServiceContext, ctx context.Context, videoId int64) (int64, error) {
	var err error
	var resp *video.CommentCountResp
	if resp, err = c.VideoRpc.CommentCount(ctx, &video.CommentCountReq{
		VideoId: videoId,
	}); err != nil {
		return 0, err
	}
	return resp.CommentCount, nil
}

func GetUserInfo(c *svc.ServiceContext, ctx context.Context, userId int64, targetId int64) (types.User, error) {
	var err error
	var resp *user.UserResp
	if resp, err = c.UserRpc.User(ctx, &user.UserReq{
		UserId:   userId,
		TargetId: targetId,
	}); err != nil {
		return types.User{}, err
	}

	return types.User{
		Id:            targetId,
		Name:          resp.Name,
		FollowCount:   *resp.FollowCount,
		FollowerCount: *resp.FollowerCount,
		IsFollow:      resp.IsFollow,
	}, nil
}
