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
	var countFans int64
	var countFollow int64
	if err := l.svcCtx.Db.Model(&models.Fans{}).
		Where(&models.Fans{UserId: in.UserId, TargetId: in.TargetId}).
		Count(&countFans).Error; err != nil {

		logx.WithContext(l.ctx).Errorf("get is friend failed: %v", err)
		return &user.IsFriendResp{
			IsFriend: false,
		}, errorx.ErrDatabaseError
	}
	if err := l.svcCtx.Db.Model(&models.Follow{}).
		Where(&models.Follow{UserId: in.UserId, TargetId: in.TargetId}).
		Count(&countFollow).Error; err != nil {

		logx.WithContext(l.ctx).Errorf("get is friend failed: %v", err)
		return &user.IsFriendResp{
			IsFriend: false,
		}, errorx.ErrDatabaseError
	}

	if countFans > 0 && countFollow > 0 {
		return &user.IsFriendResp{
			IsFriend: true,
		}, nil
	}

	logx.WithContext(l.ctx).Debugf("get is friend: %v <=> %v --> count: %v %v", in.UserId, in.TargetId, countFans, countFollow)

	return &user.IsFriendResp{
		IsFriend: false,
	}, nil
}
