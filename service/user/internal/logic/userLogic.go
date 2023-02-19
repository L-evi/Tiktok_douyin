package logic

import (
	"context"
	"train-tiktok/common/errorx"
	tool2 "train-tiktok/common/tool"
	"train-tiktok/service/identity/types/identity"
	"train-tiktok/service/user/common/tool"
	videoModels "train-tiktok/service/video/models"

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

	// check if userExist
	if exists, err := tool2.CheckUserExist(l.ctx, l.svcCtx.IdentityRpc, in.UserId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query user: %v", err)

		return nil, errorx.ErrDatabaseError
	} else if !exists {
		return nil, errorx.ErrUserNotFound
	}

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

	// get user information from database
	var userInformation models.UserInformation
	if err := l.svcCtx.Db.Model(&models.UserInformation{}).Where(&models.UserInformation{UserId: in.UserId}).Find(&userInformation).Error; err != nil {
		logx.Errorf("get user information failed: %v", err)

		return &user.UserResp{}, errorx.ErrDatabaseError
	}

	// TODO: get total favorite count
	var totalFavoriteCount int64
	// get work count
	var workCount int64
	if err := l.svcCtx.Db.Model(&videoModels.Video{}).
		Where(&videoModels.Video{UserID: in.UserId}).
		Count(&workCount).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query work count: %v", err)

		return &user.UserResp{}, errorx.ErrDatabaseError
	}

	// TODO: get favorite count
	var favoriteCount int64
	return &user.UserResp{
		Name:            rpcResp.Nickname,
		FollowCount:     &followCount,
		FollowerCount:   &followerCount,
		IsFollow:        isFollowed,
		UserId:          in.UserId,
		Avatar:          &userInformation.Avatar,
		BackgroundImage: &userInformation.BackgroundImage,
		Signature:       &userInformation.Signature,
		TotalFavorite:   &totalFavoriteCount,
		WorkCount:       &workCount,
		FavoriteCount:   &favoriteCount,
	}, nil
}
