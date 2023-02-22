package logic

import (
	"context"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByHashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByHashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByHashLogic {
	return &GetVideoByHashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByHashLogic) GetVideoByHash(in *video.GetVideoByHashReq) (*video.GetVideoByHashResp, error) {
	var myVideo models.Video
	if err := l.svcCtx.Db.Model(&models.Video{}).Where(&models.Video{Hash: in.Hash}).First(&myVideo).Error; err != nil {
		return &video.GetVideoByHashResp{
			Exists: false,
		}, nil
	}

	return &video.GetVideoByHashResp{
		Exists: true,
		Video: &video.VideoX{
			Id:       myVideo.ID,
			UserId:   myVideo.UserID,
			PlayUrl:  myVideo.PlayUrl,
			CoverUrl: myVideo.CoverUrl,
			Title:    myVideo.Title,
		},
		Position: myVideo.Position,
	}, nil
}
