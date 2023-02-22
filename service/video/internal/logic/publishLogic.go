package logic

import (
	"context"
	"strings"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/models"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {

	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *video.PublishReq) (*video.PublishResp, error) {

	// 去除开头的 "./"
	in.CoverPath = strings.TrimLeft(in.CoverPath, "./")
	in.FilePath = strings.TrimLeft(in.FilePath, "./")

	// insert to db
	if err := l.svcCtx.Db.Model(&models.Video{}).Create(&models.Video{
		UserID:   in.UserId,
		Title:    in.Title,
		PlayUrl:  in.FilePath,
		CoverUrl: in.CoverPath,
		Hash:     in.Hash,
		Position: "local",
	}).Error; err != nil {
		logx.Errorf("insert video failed: %v", err)

		return &video.PublishResp{
			Success: false,
		}, errorx.ErrDatabaseError
	}

	return &video.PublishResp{
		Success: true,
	}, nil
}
