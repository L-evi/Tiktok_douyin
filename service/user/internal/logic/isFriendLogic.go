package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/user/models"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFriendLogic) IsFriend(in *user.IsFriendReq) (*user.IsFriendResp, error) {
	if in.UserId == in.TargetId {
		return &user.IsFriendResp{
			IsFriend: false,
		}, nil
	}
	if in.UserId == 0 || in.TargetId == 0 {
		return &user.IsFriendResp{
			IsFriend: false,
		}, nil
	}
	var count int64
	if err := l.svcCtx.Db.Model(&models.Fans{}).
		Where(&models.Fans{UserId: in.UserId, TargetId: in.TargetId}).
		Where(&models.Fans{UserId: in.TargetId, TargetId: in.UserId}).
		Count(&count).Error; err != nil {

		logx.WithContext(l.ctx).Errorf("get is friend failed: %v", err)
		return &user.IsFriendResp{
			IsFriend: false,
		}, errorx.ErrDatabaseError
	}

	if count > 0 {
		return &user.IsFriendResp{
			IsFriend: true,
		}, nil
	}

	return &user.IsFriendResp{
		IsFriend: false,
	}, nil
}
