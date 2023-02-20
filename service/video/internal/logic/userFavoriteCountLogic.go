package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"train-tiktok/service/video/common/rediskeyutil"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFavoriteCountLogic {
	return &UserFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserFavoriteCountLogic) UserFavoriteCount(in *video.UserFavoriteCountReq) (*video.UserFavoriteCountResp, error) {
	userId := in.UserId
	_redisKey := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix).GetUserKey(userId)

	count, err := l.svcCtx.Rdb.ZCard(l.ctx, _redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return &video.UserFavoriteCountResp{
			FavoriteCount: 0,
		}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("redis ZCard error: %v", err)

		return &video.UserFavoriteCountResp{}, err
	}

	return &video.UserFavoriteCountResp{
		FavoriteCount: count,
	}, nil
}
