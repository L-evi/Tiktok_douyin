package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/types/identity"
	"train-tiktok/service/user/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/models"
	"train-tiktok/service/user/types/user"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// User get user info
//
//	对应 /douyin/user 接口 用于获取用户相关信息
func (l *UserLogic) User(in *user.UserReq) (*user.UserResp, error) {
	var followerCount int64
	var followCount int64
	var isFollowed bool

	var err error

	// get followCount
	if err = l.svcCtx.Db.Model(&models.Follow{}).
		Where(&models.Follow{UserId: in.TargetId}).
		Count(&followCount).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query followCount: %v", err)
		return &user.UserResp{}, errorx.ErrDatabaseError

	}

	// get followerCount
	if err = l.svcCtx.Db.Model(&models.Fans{}).
		Where(&models.Fans{UserId: in.TargetId}).
		Count(&followerCount).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query followerCount: %v", err)

		return &user.UserResp{}, errorx.ErrDatabaseError
	}

	if isFollowed, err = tool.IsFollowing(l.ctx, l.svcCtx.Db, in.UserId, in.TargetId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query isFollowed: %v", err)

		return &user.UserResp{}, errorx.ErrDatabaseError
	}

	// Get Username From Identity
	var rpcResp *identity.GetUserInfoResp
	if rpcResp, err = l.svcCtx.IdentityRpc.GetUserInfo(l.ctx, &identity.GetUserInfoReq{
		UserId: in.TargetId,
	}); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query username: %v", err)

		return &user.UserResp{}, errorx.ErrDatabaseError
	}

	return &user.UserResp{
		Name:          rpcResp.Nickname,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollowed,
	}, nil
}
