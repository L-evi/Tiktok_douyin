package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/rediskeyutil"
	"train-tiktok/service/video/common/tool"
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

// FavoriteCount 获取某视频点赞数
func (l *FavoriteCountLogic) FavoriteCount(in *video.FavoriteCountReq) (*video.FavoriteCountResp, error) {

	// check video exists
	if exists, err := tool.CheckVideoExists(l.svcCtx.Db, in.VideoId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query video, err: %v", err)

		return nil, err
	} else if !exists {
		return nil, errx.ErrVideoNotFound
	}

	// connect to redis
	rdb := l.svcCtx.Rdb

	// get favorite count
	_redisKey := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix).GetVideoFavoriteKey(in.VideoId)

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
