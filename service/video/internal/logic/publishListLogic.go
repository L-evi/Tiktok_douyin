package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/common/tool"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishListLogic) PublishList(in *video.PublishListReq) (*video.PublishListResp, error) {
	// query
	var videos []models.Video
	if err := l.svcCtx.Db.Model(&models.Video{}).Where(&models.Video{UserID: in.UserId}).
		Order("create_at desc").Find(&videos).
		Error; errors.Is(err, gorm.ErrRecordNotFound) {

		return &video.PublishListResp{}, errx.ErrNoLatestVideo
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("PublishList Feed sql err: %v", err)

		return &video.PublishListResp{}, err
	}

	// 处理数据
	var videoList []*video.VideoX
	videoList = make([]*video.VideoX, 0, len(videos))

	for _, v := range videos {
		v.PlayUrl, v.CoverUrl = tool.HandleVideoUrl(l.svcCtx, v.Position, v.PlayUrl, v.CoverUrl)

		// insert videoList
		videoList = append(videoList, &video.VideoX{
			Id:       v.ID,
			UserId:   v.UserID,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		})
	}

	return &video.PublishListResp{
		VideoList: videoList,
	}, nil
}
