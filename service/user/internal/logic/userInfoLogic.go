package logic

import (
	"context"
	"train-tiktok/service/user/models"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoReq) (*user.UserInfoResp, error) {
	var userInfo models.UserInfo
	if err := l.svcCtx.Db.Model(&models.UserInfo{}).
		Where(&models.UserInfo{UserId: in.UserId}).
		Find(&userInfo).Error; err != nil {
		logx.Errorf("get user infomation failed: %v", err)

		return &user.UserInfoResp{}, err
	}

	return &user.UserInfoResp{
		UserInfo: &user.UserInfo{
			UserId:          userInfo.UserId,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
			TotalFavorite:   userInfo.TotalFavorite,
			WorkCount:       userInfo.WorkCount,
			FavoriteCount:   userInfo.FavoriteCount,
		},
	}, nil
}
