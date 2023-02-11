package logic

import (
	"context"
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
	if in.ActionType == 1 {
		// add comment
		var Comment = &models.Comment{
			VideoID: in.VideoId,
			UserID:  l.ctx.Value("user_id").(int64),
			Content: in.CommentText,
		}
		if res := l.svcCtx.Db.Create(&Comment); res.Error != nil || res.RowsAffected == 0 {
			logx.Errorf("failed to create comment, err: %v", res.Error)

			return &video.CommentActionResp{}, errorx.ErrSystemError
		}

		return &video.CommentActionResp{
			Comment: &video.Comment{
				Id:         Comment.ID,
				Content:    Comment.Content,
				CreateDate: Comment.CreateDate,
				UserId:     l.ctx.Value("user_id").(int64),
			},
		}, nil
	} else if in.ActionType == 2 {
		// delete comment
		if res := l.svcCtx.Db.Delete(&models.Comment{}, in.CommentId); res.Error != nil || res.RowsAffected == 0 {
			logx.Errorf("failed to delete comment, err: %v", res.Error)

			return &video.CommentActionResp{}, errorx.ErrSystemError
		}

		return &video.CommentActionResp{}, nil
	}

	return &video.CommentActionResp{}, nil
}
