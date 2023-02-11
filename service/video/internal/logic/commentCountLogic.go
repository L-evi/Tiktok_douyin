package logic

import (
	"context"
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
	// get comment count
	// TODO  to redis

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
