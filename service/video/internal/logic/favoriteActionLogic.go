package logic

import (
	"context"
	"fmt"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteActionLogic) FavoriteAction(in *video.FavoriteActionReq) (*video.FavoriteActionResp, error) {
	// connect to redis
	rdb := l.svcCtx.Rdb

	_redisKey := fmt.Sprintf("%s:favorite:count:%d", l.svcCtx.Config.Redis.Prefix, in.VideoId)

	// favorite action
	switch in.ActionType {
	case 1:
		if _, err := rdb.Incr(l.ctx, _redisKey).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Incr error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		break
	case 2:
		if _, err := rdb.Decr(l.ctx, _redisKey).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Decr error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		break
	default:
		return &video.FavoriteActionResp{}, errorx.ErrSystemError
	}

	return &video.FavoriteActionResp{}, nil
}
