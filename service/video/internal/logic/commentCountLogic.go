package logic

import (
	"context"
	"strconv"
	"train-tiktok/service/video/models"

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
	// connect to redis
	rdb := l.svcCtx.Rdb
	var ctx context.Context
	// get comment count from redis
	key := "tiktok:comment:count:" + strconv.FormatInt(in.VideoId, 10)
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		logx.Errorf("redis Get error: %v", err)

		return &video.CommentCountResp{}, err
	}
	if result == "" {
		// get comment count from database
		var commentCount int64
		if err := l.svcCtx.Db.Model(&models.Comment{}).
			Where(&models.Comment{ID: in.VideoId}).
			Count(&commentCount).Error; err != nil {

			return &video.CommentCountResp{}, err
		}

		return &video.CommentCountResp{
			CommentCount: commentCount,
		}, nil
	}
	commentCount, err := strconv.Atoi(result)
	if err != nil {
		logx.Errorf("string to int error: %v", err)

		return &video.CommentCountResp{}, err
	}

	return &video.CommentCountResp{
		CommentCount: (int64)(commentCount),
	}, nil
}
