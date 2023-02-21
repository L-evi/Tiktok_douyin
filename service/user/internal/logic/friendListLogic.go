package logic

import (
	"context"
	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/models"
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

	// models.Fans(targetId, userId) targetId 新增 userId 为粉丝 // uid 是 tid 的粉丝
	// models.Follow(userId, targetId) userId 新增 targetId 为关注对象 // uid 关注 tid
	// 求 user_id = 1 的好友列表 (互关)
	//	// X select target_id from fans where user_id = 1 and target_id in (select user_id from follows where target_id = 1)
	// select user_id from follows where target_id = 2 and user_id in (select user_id from fans where target_id = 2);
	result := l.svcCtx.Db.Model(&models.Follow{}).Select([]string{"user_id"}).Where("target_id = ?", in.UserId).
		Where("user_id in (?)", l.svcCtx.Db.Model(&models.Fans{}).Select("user_id").Where("target_id = ?", in.UserId)).
		Pluck("user_id", &userId)
	if result.Error != nil {
		logx.WithContext(l.ctx).Error("get friend list failed: %v", result.Error)

		return &user.FriendListResp{}, nil
	}

	logx.WithContext(l.ctx).Debugf("get friend list success: %v", userId)

	return &user.FriendListResp{
		UserIdList: userId,
	}, nil
}
