package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/tool"
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

	// check video exists
	if exists, err := tool.CheckVideoExists(l.svcCtx.Db, in.VideoId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query video, err: %v", err)

		return nil, err
	} else if !exists {
		return nil, errx.ErrVideoNotFound
	}

	// connect to redis
	rdb := l.svcCtx.Rdb

	// 记录视频点赞数
	_countKey := fmt.Sprintf("%s:favorite_count:%d", l.svcCtx.Config.RedisConf.Prefix, in.VideoId)
	// 记录用户是否点赞该视频
	_userKey := fmt.Sprintf("%s:favorite_user:%d", l.svcCtx.Config.RedisConf.Prefix, in.UserId)

	videoIdStr := strconv.FormatInt(in.VideoId, 10)

	timeNow := time.Now().Unix()

	// favorite action
	switch in.ActionType {
	case 1:
		// check if already favorite
		// 事物
		tran := rdb.Watch(l.ctx, func(tx *redis.Tx) error {
			// 检查用户是否已经点赞过该视频
			if exists, err := tx.ZScore(l.ctx, _userKey, videoIdStr).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("redis ZScore error: %v", err)
				return err
			} else if exists != 0 {
				return errx.ErrAlreadyFavorite
			}
			// 事物
			pipe := tx.Pipeline()
			if _, err := pipe.Incr(l.ctx, _countKey).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Incr error: %v", err)
				return err
			}
			if _, err := pipe.ZAdd(l.ctx, _userKey, redis.Z{
				Score:  float64(timeNow),
				Member: videoIdStr,
			}).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("redis ZADD error: %v", err)
				return err
			}
			if _, err := pipe.Exec(l.ctx); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
				return err
			}
			return nil
		}, _countKey, _userKey)
		if tran != nil {
			return &video.FavoriteActionResp{}, tran
		}
		break
	case 2:
		tran := rdb.Watch(l.ctx, func(tx *redis.Tx) error {
			pipe := tx.Pipeline()
			if _, err := pipe.Decr(l.ctx, _countKey).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Decr error: %v", err)
				return err
			}
			if _, err := pipe.ZRem(l.ctx, _userKey, videoIdStr).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("redis ZRem error: %v", err)
				return err
			}
			if _, err := pipe.Exec(l.ctx); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
				return err
			}
			return nil
		}, _countKey, _userKey)
		if tran != nil {
			return &video.FavoriteActionResp{}, tran
		}
		break
	default:
		return &video.FavoriteActionResp{}, errorx.ErrInvalidParameter
	}

	return &video.FavoriteActionResp{}, nil
}
