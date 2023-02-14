package logic

import (
	"context"
	"time"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/tool"
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

	// check video exists
	if exists, err := tool.CheckVideoExists(l.svcCtx.Db, in.VideoId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query video, err: %v", err)

		return nil, err
	} else if !exists {
		return nil, errx.ErrVideoNotFound
	}

	var commentList []models.Comment
	res := l.svcCtx.Db.Model(&video.Comment{}).Where(&models.Comment{VideoID: in.VideoId}).
		Order("create_at desc").Find(&commentList)
	if res.Error != nil {
		logx.Errorf("get comment list failed: %v", res.Error)

		return &video.CommentListResp{}, res.Error
	}

	var list []*video.Comment
	list = make([]*video.Comment, 0, len(commentList))
	for _, v := range commentList {
		CreateDate := time.Unix(v.CreateAt, 0).Format("01-02")
		list = append(list, &video.Comment{
			Id:         v.ID,
			UserId:     v.UserID,
			Content:    v.Content,
			CreateDate: CreateDate,
		})
	}
	//log.Println(list)

	return &video.CommentListResp{
		CommentList: list,
	}, nil
}
