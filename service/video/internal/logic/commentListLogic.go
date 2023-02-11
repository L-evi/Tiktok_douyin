package logic

import (
	"context"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {

	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListLogic) CommentList(in *video.CommentListReq) (*video.CommentListResp, error) {
	var commentList []models.Comment
	res := l.svcCtx.Db.Where(&video.Comment{Id: in.VideoId}).Find(&commentList)
	if res.Error != nil {
		logx.Errorf("get comment list failed: %v", res.Error)

		return &video.CommentListResp{}, res.Error
	}

	var list []*video.Comment
	for _, v := range commentList {
		list = append(list, &video.Comment{
			Id:         v.ID,
			UserId:     v.UserID,
			Content:    v.Content,
			CreateDate: v.CreateDate,
		})
	}

	return &video.CommentListResp{
		CommentList: list,
	}, nil
}
