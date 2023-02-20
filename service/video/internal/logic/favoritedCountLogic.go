package logic

import (
	"context"
	"train-tiktok/service/video/common/rediskeyutil"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoritedCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoritedCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoritedCountLogic {
	return &FavoritedCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoritedCountLogic) FavoritedCount(in *video.FavoritedCountReq) (*video.FavoritedCountResp, error) {
	userId := in.UserId
	_redisKey := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix).GetPublisherFavoriteKey(userId)

	count, err := l.svcCtx.Rdb.Get(l.ctx, _redisKey).Int64()
	if err != nil {
		logx.WithContext(l.ctx).Errorf("redis ZCard error: %v", err)

		return &video.FavoritedCountResp{}, err
	}

	return &video.FavoritedCountResp{
		FavoriteCount: count,
	}, nil
}
