package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/tool"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentCountLogic {

	return &CommentCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentCountLogic) CommentCount(in *video.CommentCountReq) (*video.CommentCountResp, error) {

	// check video exists
	if exists, err := tool.CheckVideoExists(l.svcCtx.Db, in.VideoId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query video, err: %v", err)

		return nil, err
	} else if !exists {
		return nil, errx.ErrVideoNotFound
	}

	// connect to redis
	rdb := l.svcCtx.Rdb

	// get comment count from redis
	_redisKey := fmt.Sprintf("%s:comment_count:%d", l.svcCtx.Config.RedisConf.Prefix, in.VideoId)
	var result int64
	var err error
	if result, err = rdb.Get(l.ctx, _redisKey).Int64(); !errors.Is(err, redis.Nil) && err != nil {
		logx.Errorf("redis Get error: %v", err)

		return &video.CommentCountResp{}, err
	}

	return &video.CommentCountResp{
		CommentCount: result,
	}, nil
}
