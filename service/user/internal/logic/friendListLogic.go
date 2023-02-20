package logic

import (
	"context"
	"train-tiktok/service/user/models"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *user.FriendListReq) (*user.FriendListResp, error) {
	// get friends from fans join follow
	var userId []int64
	result := l.svcCtx.Db.Joins("follow").
		Where(&models.Follow{UserId: in.UserId}).
		Where(&models.Fans{UserId: in.UserId}).
		Pluck("target_id", &userId)
	if result.Error != nil {
		logx.Error("get friend list failed: %v", result.Error)

		return &user.FriendListResp{}, nil
	}

	return &user.FriendListResp{
		UserIdList: userId,
	}, nil
}
