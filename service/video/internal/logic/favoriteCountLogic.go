package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteCountLogic {
	return &FavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteCountLogic) FavoriteCount(in *video.FavoriteCountReq) (*video.FavoriteCountResp, error) {
	// connect to redis
	rdb := l.svcCtx.Rdb

	// get favorite count
	_redisKey := fmt.Sprintf("%s:favorite_count:%d", l.svcCtx.Config.RedisConf.Prefix, in.VideoId)

	var err error
	var result int64
	if result, err = rdb.Get(l.ctx, _redisKey).Int64(); !errors.Is(err, redis.Nil) && err != nil {
		logx.WithContext(l.ctx).Errorf("redis Get error: %v", err)
		return &video.FavoriteCountResp{}, err
	}

	return &video.FavoriteCountResp{
		FavoriteCount: result,
	}, nil
}
