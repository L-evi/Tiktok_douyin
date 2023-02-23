package logic

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/rediskeyutil"
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
	var err error
	var videoPublisherId int64 // 视频发布者的 userId

	if videoPublisherId, err = tool.GetVideoUserId(l.svcCtx.Db, in.VideoId); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errx.ErrVideoNotFound
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query video, err: %v", err)

		return nil, err
	}

	// connect to redis
	rdb := l.svcCtx.Rdb

	keys := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix)
	// 记录视频点赞数
	_videoCountKey := keys.GetVideoFavoriteKey(in.VideoId)
	// 记录用户是否点赞该视频
	_userKey := keys.GetUserKey(in.UserId)
	// 记录视频发布用户的 获赞数
	_favoritedKey := keys.GetPublisherFavoriteKey(videoPublisherId)

	videoIdStr := strconv.FormatInt(in.VideoId, 10)

	timeNow := time.Now().Unix()

	// favorite action
	switch in.ActionType {
	case 1:
		// check if already favorite
		// 事物
		tran := rdb.Watch(l.ctx, func(tx *redis.Tx) error {
			// 检查用户是否已经点赞过该视频
			if exists, err := tx.ZScore(l.ctx, _userKey, videoIdStr).Result(); !errors.Is(err, redis.Nil) && err != nil {
				logx.WithContext(l.ctx).Errorf("redis ZScore error: %v", err)

				return err
			} else if exists != 0 {
				// 已经点赞了，不需要再点赞 向前端返回点赞成功避免错误
				return nil
				// return errx.ErrAlreadyFavorite
			}

			// 事物
			pipe := tx.Pipeline()
			{
				if err := pipe.Incr(l.ctx, _videoCountKey).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis Incr error: %v", err)
					return err
				}
				if err := pipe.ZAdd(l.ctx, _userKey, redis.Z{
					Score:  float64(timeNow),
					Member: videoIdStr,
				}).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis ZADD error: %v", err)
					return err
				}
				if err := pipe.Incr(l.ctx, _favoritedKey).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis incr error: %v", err)

					return err
				}
			}
			if _, err := pipe.Exec(l.ctx); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
				return err
			}
			return nil
		}, _videoCountKey, _favoritedKey, _userKey)
		if tran != nil {
			return &video.FavoriteActionResp{}, tran
		}
		break
	case 2:
		tran := rdb.Watch(l.ctx, func(tx *redis.Tx) error {
			pipe := tx.Pipeline()
			{
				if err := pipe.Decr(l.ctx, _videoCountKey).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis Decr error: %v", err)
					return err
				}
				if err := pipe.ZRem(l.ctx, _userKey, videoIdStr).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis ZRem error: %v", err)
					return err
				}
				if err := pipe.Decr(l.ctx, _favoritedKey).Err(); err != nil {
					logx.WithContext(l.ctx).Errorf("redis incr error: %v", err)

					return err
				}
			}

			if _, err := pipe.Exec(l.ctx); err != nil {
				logx.WithContext(l.ctx).Errorf("redis Exec error: %v", err)
				return err
			}
			return nil
		}, _videoCountKey, _favoritedKey, _userKey)
		if tran != nil {
			return &video.FavoriteActionResp{}, tran
		}
		break
	default:
		return &video.FavoriteActionResp{}, errorx.ErrInvalidParameter
	}

	return &video.FavoriteActionResp{}, nil
}
