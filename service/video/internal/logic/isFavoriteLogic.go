package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
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

	_userKey := fmt.Sprintf("%s:favorite_user:%d", l.svcCtx.Config.RedisConf.Prefix, in.UserId)
	videoIdStr := strconv.FormatInt(in.VideoId, 10)

	var err error
	var isFavorite bool
	if isFavorite, err = rdb.HExists(l.ctx, _userKey, videoIdStr).Result(); errors.Is(err, redis.Nil) {
		return &video.IsFavoriteResp{
			IsFavorite: false,
		}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("redis HExists error: %v", err)
		return &video.IsFavoriteResp{}, err
	}

	return &video.IsFavoriteResp{
		IsFavorite: isFavorite,
	}, nil
}
