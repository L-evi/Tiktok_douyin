package logic

import (
	"context"
	"log"
	"strconv"
	"train-tiktok/common/redisutil"

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
	// todo: add your logic here and delete this line
	// connect to redis
	rdb := redisutil.New(redisutil.RedisConf{
		Addr:        l.svcCtx.Config.Redis.Addr,
		Password:    l.svcCtx.Config.Redis.Password,
		DB:          l.svcCtx.Config.Redis.DB,
		MinIdle:     l.svcCtx.Config.Redis.MinIdle,
		PoolSize:    l.svcCtx.Config.Redis.PoolSize,
		MaxLifeTime: l.svcCtx.Config.Redis.MaxLifeTime,
	})
	var ctx context.Context
	defer rdb.Close()
	// favorite action
	if in.ActionType == 1 {
		// get favorite count
		key := strconv.FormatInt(in.VideoId, 10) + "_favorite_count"
		result, err := rdb.Get(ctx, key).Result()
		if err != nil {
			log.Printf("redis Get error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		if result == "" {
			// new key and set value
			err = rdb.Set(ctx, key, 1, 0).Err()
			if err != nil {
				log.Printf("redis Set error: %v", err)
				return &video.FavoriteActionResp{}, err
			}
		} else {
			// get value and add 1
			count, err := strconv.Atoi(result)
			if err != nil {
				log.Printf("string to int error: %v", err)
				return &video.FavoriteActionResp{}, err
			}
			count++
			err = rdb.Set(ctx, key, count, 0).Err()
			if err != nil {
				log.Printf("redis Set error: %v", err)
				return &video.FavoriteActionResp{}, err
			}
		}
	}
	return &video.FavoriteActionResp{}, nil
}
