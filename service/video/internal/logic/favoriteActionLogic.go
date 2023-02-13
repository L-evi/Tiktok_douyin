package logic

import (
	"context"
	"fmt"
	"strconv"
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

	// 记录视频点赞数
	_countKey := fmt.Sprintf("%s:favorite_count:%d", l.svcCtx.Config.RedisConf.Prefix, in.VideoId)
	// 记录用户是否点赞该视频
	_userKey := fmt.Sprintf("%s:favorite_user:%d", l.svcCtx.Config.RedisConf.Prefix, in.UserId)

	videoIdStr := strconv.FormatInt(in.VideoId, 10)

	// favorite action
	switch in.ActionType {
	case 1:
		pipe := rdb.Pipeline()
		if _, err := pipe.Incr(l.ctx, _countKey).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Incr error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		if _, err := pipe.HSet(l.ctx, _userKey, videoIdStr, 1).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis HSet error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		if _, err := pipe.Exec(l.ctx); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		break
	case 2:
		pipe := rdb.Pipeline()
		if _, err := pipe.Decr(l.ctx, _countKey).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Incr error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		if _, err := pipe.HDel(l.ctx, _userKey, videoIdStr).Result(); err != nil {
			logx.WithContext(l.ctx).Errorf("redis HSet error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		if _, err := pipe.Exec(l.ctx); err != nil {
			logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		break
	default:
		return &video.FavoriteActionResp{}, errorx.ErrSystemError
	}

	return &video.FavoriteActionResp{}, nil
}
