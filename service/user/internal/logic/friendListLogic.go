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

	// x SELECT f.target_id FROM fans f LEFT JOIN follows fo ON f.target_id = fo.user_id WHERE f.user_id = 1 AND fo.target_id = 1;
	// select fans.user_id from fans where fans.target_id = 1 UNION select follows.target_id from follows where follows.user_id = 1;
	result := l.svcCtx.Db.Raw("? UNION ?",
		l.svcCtx.Db.Table("fans").Select("fans.user_id").Where("fans.target_id = ?", in.UserId),
		l.svcCtx.Db.Table("follows").Select("follows.target_id").Where("follows.user_id = ?", in.UserId),
	).Scan(&userId)

	if result.Error != nil {
		logx.WithContext(l.ctx).Error("get friend list failed: %v", result.Error)

		return &user.FriendListResp{}, nil
	}

	logx.WithContext(l.ctx).Debugf("get friend list success: %v", userId)

	return &user.FriendListResp{
		UserIdList: userId,
	}, nil
}
