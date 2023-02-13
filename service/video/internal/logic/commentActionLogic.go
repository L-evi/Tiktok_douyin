package logic

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/models"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {

	return &CommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentActionLogic) CommentAction(in *video.CommentActionReq) (*video.CommentActionResp, error) {
	rdb := l.svcCtx.Rdb

	// get comment count from redis
	_redisKey := fmt.Sprintf("%s:comment_count:%d", l.svcCtx.Config.RedisConf.Prefix, in.VideoId)

	if in.ActionType == 1 {
		// add comment
		var Comment = &models.Comment{
			VideoID: in.VideoId,
			UserID:  in.UserId,
			Content: in.CommentText,
		}

		if err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
			if res := l.svcCtx.Db.Create(&Comment); res.Error != nil || res.RowsAffected == 0 {
				logx.WithContext(l.ctx).Errorf("failed to create comment, err: %v", res.Error)
				return errorx.ErrSystemError
			}
			if _, err := rdb.Incr(l.ctx, _redisKey).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("failed to incr comment count, err: %v", err)
				return errorx.ErrSystemError
			}
			return nil
		}); err != nil {
			return &video.CommentActionResp{}, err
		}

		// 秒级时间戳 转换为 mm-dd
		CreateDate := time.Unix(Comment.CreateAt, 0).Format("01-02")

		return &video.CommentActionResp{
			Comment: &video.Comment{
				Id:         Comment.ID,
				Content:    Comment.Content,
				CreateDate: CreateDate,
				UserId:     in.UserId,
			},
		}, nil
	} else if in.ActionType == 2 {
		// delete comment
		if err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
			if res := l.svcCtx.Db.Delete(&models.Comment{}, in.CommentId); res.Error != nil || res.RowsAffected == 0 {
				logx.Errorf("failed to delete comment, err: %v", res.Error)

				return errorx.ErrSystemError
			}
			if _, err := rdb.Decr(l.ctx, _redisKey).Result(); err != nil {
				logx.WithContext(l.ctx).Errorf("failed to decr comment count, err: %v", err)
				return errorx.ErrSystemError
			}
			return nil
		}); err != nil {
			return &video.CommentActionResp{}, err
		}

		return &video.CommentActionResp{}, nil
	}

	return &video.CommentActionResp{}, nil
}
