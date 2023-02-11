package logic

import (
	"context"
	"log"
	"strconv"
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
	var ctx context.Context
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
	} else if in.ActionType == 2 {
		// conceal favorite
		key := strconv.FormatInt(in.VideoId, 10) + "_favorite_count"
		result, err := rdb.Get(ctx, key).Result()
		if err != nil {
			log.Printf("redis Get error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		count, err := strconv.Atoi(result)
		if err != nil {
			log.Printf("string to int error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
		count--
		err = rdb.Set(ctx, key, count, 0).Err()
		if err != nil {
			log.Printf("redis Set error: %v", err)
			return &video.FavoriteActionResp{}, err
		}
	}

	return &video.FavoriteActionResp{}, nil
}
