package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
	"train-tiktok/service/user/models"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *user.FollowListReq) (*user.FollowListResp, error) {
	if in.UserId == 0 {
		return nil, errorx.ErrInvalidParameter
	}
	var users []int64
	if err := l.svcCtx.Db.Model(&models.Follow{}).Where(models.Fans{
		UserId: in.UserId,
	}).Pluck("target_id", &users).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return &user.FollowListResp{}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query follows: %v", err)

		return nil, errorx.ErrDatabaseError
	}

	return &user.FollowListResp{
		UserIds: users,
	}, nil
}
