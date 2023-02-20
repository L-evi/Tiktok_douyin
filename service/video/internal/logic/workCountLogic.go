package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWorkCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkCountLogic {
	return &WorkCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WorkCountLogic) WorkCount(in *video.WorkCountReq) (*video.WorkCountResp, error) {
	var count int64
	var err error
	if err = l.svcCtx.Db.Model(&models.Video{}).Where(models.Video{UserID: in.UserId}).Count(&count).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("db count error: %v", err)

		return &video.WorkCountResp{}, errorx.ErrDatabaseError
	}
	return &video.WorkCountResp{
		WorkCount: count,
	}, nil
}
