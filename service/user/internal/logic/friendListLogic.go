package logic

import (
	"context"
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

	//SELECT f.target_id FROM fans f LEFT JOIN follows fo ON f.target_id = fo.user_id WHERE f.user_id = 1 AND fo.target_id = 1;
	result := l.svcCtx.Db.Table("fans").Select("fans.target_id").
		Joins("LEFT JOIN follows ON fans.target_id = follows.user_id").
		Where("fans.user_id = ? AND follows.target_id = ?", in.UserId, in.UserId).
		Find(&userId)
	if result.Error != nil {
		logx.Error("get friend list failed: %v", result.Error)

		return &user.FriendListResp{}, nil
	}

	return &user.FriendListResp{
		UserIdList: userId,
	}, nil
}
