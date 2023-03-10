package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/user/models"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FollowerList 获取某用户的粉丝
func (l *FollowerListLogic) FollowerList(in *user.FollowerListReq) (*user.FollowerListResp, error) {
	if in.UserId == 0 {
		return nil, errorx.ErrInvalidParameter
	}
	// checkIf targetId exists

	var users []int64
	if err := l.svcCtx.Db.Model(&models.Fans{}).Where(models.Fans{
		UserId: in.UserId,
	}).Pluck("target_id", &users).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query fans: %v", err)

		return nil, errorx.ErrDatabaseError
	}

	return &user.FollowerListResp{
		UserIds: users,
	}, nil
}
