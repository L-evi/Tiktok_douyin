package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
	"train-tiktok/service/identity/common/errutil"
	"train-tiktok/service/identity/common/userutil"
	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/models"
	"train-tiktok/service/identity/types/identity"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *identity.LoginReq) (*identity.LoginResp, error) {

	if err := userutil.VerifyUsername(in.Username); err != nil {
		return &identity.LoginResp{}, err
	}

	// check User info
	var User models.User
	if err := l.svcCtx.Db.Select("id, username, password").
		Where(&models.User{Username: in.Username}).
		First(&User).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return &identity.LoginResp{}, errutil.ErrWrongIdentity
	} else if err != nil {
		logx.Errorf("failed to query user: %v", err)
		return &identity.LoginResp{}, errorx.ErrDatabaseError
	}

	// check pwd
	if ok, err := tool.VerifyPassword(in.Password, User.Password); err != nil || !ok {
		return &identity.LoginResp{}, errutil.ErrWrongIdentity
	}

	// create token
	token, err := userutil.GenerateJwt(l.svcCtx, User.ID, in.Username)
	if err != nil {
		logx.Errorf("failed to generate jwt: %v", err)
		return &identity.LoginResp{}, errorx.ErrSystemError
	}

	return &identity.LoginResp{
		Response: &identity.Resp{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: User.ID,
		Token:  token,
	}, nil
}
