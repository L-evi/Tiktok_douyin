package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/models"
	"train-tiktok/service/identity/types/identity"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo get user info
// return [nickname, username]
func (l *GetUserInfoLogic) GetUserInfo(in *identity.GetUserInfoReq) (*identity.GetUserInfoResp, error) {
	var userId = in.UserId
	var _user models.User
	var _userInfo models.UserInformation

	if err := l.svcCtx.Db.Model(&models.User{}).Select([]string{"username"}).
		Where(&models.User{ID: userId}).
		First(&_user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	} else if err != nil {

		logx.Errorf("failed to query user: %v", err)
		return nil, errorx.ErrDatabaseError
	}

	if err := l.svcCtx.Db.Model(&models.UserInformation{}).Select([]string{"nickname", "background_image", "avatar", "signature"}).
		Where(&models.UserInformation{UserId: userId}).
		First(&_user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	} else if err != nil {

		logx.Errorf("failed to query user: %v", err)
		return nil, errorx.ErrDatabaseError
	}

	return &identity.GetUserInfoResp{
		Username:        _user.Username,
		Nickname:        _userInfo.Nickname,
		Avatar:          _userInfo.Avatar,
		BackgroundImage: _userInfo.BackgroundImage,
		Signature:       _userInfo.Signature,
	}, nil
}
