package rpcutil

import (
	"context"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	"train-tiktok/service/identity/identityclient"
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
	// 请求 userRpc 获取用户基本信息
	if resp, err = c.UserRpc.User(ctx, &user.UserReq{
		UserId:   userId,
		TargetId: targetId,
	}); err != nil {
		return types.User{}, err
	}
	// 请求 videoRpc 获取用户作品/点赞信息
	var workRpc *video.WorkCountResp
	if workRpc, err = c.VideoRpc.WorkCount(ctx, &video.WorkCountReq{
		UserId: targetId,
	}); err != nil {
		return types.User{}, err
	}

	var favoriteRpc *video.FavoriteCountResp
	if favoriteRpc, err = c.VideoRpc.FavoriteCount(ctx, &video.FavoriteCountReq{
		VideoId: 0,
	}); err != nil {
		return types.User{}, err
	}

	var favoritedCount *video.FavoritedCountResp
	if favoritedCount, err = c.VideoRpc.FavoritedCount(ctx, &video.FavoritedCountReq{
		UserId: targetId,
	}); err != nil {
		return types.User{}, err
	}

	// 请求 identityRpc 获取信息
	var identityRpc *identityclient.GetUserInfoResp
	if identityRpc, err = c.IdentityRpc.GetUserInfo(ctx, &identityclient.GetUserInfoReq{
		UserId: targetId,
	}); err != nil {
		return types.User{}, err
	}

	return types.User{
		Id:              targetId,
		Name:            identityRpc.Nickname,
		FollowCount:     *resp.FollowCount,
		FollowerCount:   *resp.FollowerCount,
		IsFollow:        resp.IsFollow,
		Avatar:          identityRpc.Avatar,
		Signature:       identityRpc.Signature,
		BackgroundImage: identityRpc.BackgroundImage,
		WorkCount:       workRpc.WorkCount,
		TotalFavorited:  favoritedCount.FavoriteCount,
		FavoriteCount:   favoriteRpc.FavoriteCount,
	}, nil
}
