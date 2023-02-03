package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/models"

	"train-tiktok/service/identity/internal/svc"
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

	err := l.svcCtx.Db.Model(&models.User{}).Select([]string{"username", "nickname"}).
		Where(&models.User{ID: userId}).
		First(&_user).Error
	if err != nil {
		logx.Errorf("failed to query user: %v", err)
		return nil, errorx.ErrDatabaseError
	}

	return &identity.GetUserInfoResp{
		Username: _user.Username,
		Nickname: _user.Nickname,
	}, nil
}
