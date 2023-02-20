package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"train-tiktok/service/video/common/rediskeyutil"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteLogic {
	return &IsFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFavoriteLogic) IsFavorite(in *video.IsFavoriteReq) (*video.IsFavoriteResp, error) {
	rdb := l.svcCtx.Rdb
	if in.UserId == 0 || in.VideoId == 0 {
		return &video.IsFavoriteResp{
			IsFavorite: false,
		}, nil
	}

	_userKey := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix).GetUserKey(in.UserId)
	
	videoIdStr := strconv.FormatInt(in.VideoId, 10)

	var err error
	var isFavorite float64
	if isFavorite, err = rdb.ZScore(l.ctx, _userKey, videoIdStr).Result(); errors.Is(err, redis.Nil) {
		return &video.IsFavoriteResp{
			IsFavorite: false,
		}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("redis HExists error: %v", err)
		return &video.IsFavoriteResp{}, err
	}

	return &video.IsFavoriteResp{
		IsFavorite: isFavorite > 0,
	}, nil
}
