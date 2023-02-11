package logic

import (
	"context"
	"strconv"

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
	var ctx context.Context

	// get favorite count
	key := strconv.FormatInt(in.VideoId, 10) + "_favorite_count"
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		logx.Errorf("redis Get error: %v", err)
		return &video.FavoriteCountResp{}, err
	}

	if result != "" {
		// get value
		count, err := strconv.Atoi(result)
		if err != nil {
			logx.Errorf("string to int error: %v", err)

			return &video.FavoriteCountResp{}, err
		}

		return &video.FavoriteCountResp{
			FavoriteCount: int64(count),
		}, nil
	}

	return &video.FavoriteCountResp{
		FavoriteCount: 0,
	}, nil
}
