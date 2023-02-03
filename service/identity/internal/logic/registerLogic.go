package logic

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
	"train-tiktok/service/identity/common/userutil"
	"train-tiktok/service/identity/models"

	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/types/identity"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *identity.RegisterReq) (*identity.RegisterResp, error) {

	// safe check
	if err := userutil.VerifyUsername(in.Username); err != nil {
		return &identity.RegisterResp{}, err
	}

	if err := userutil.VerifyPwd(in.Password); err != nil {
		return &identity.RegisterResp{}, err
	}

	// check if username exists
	if err := userutil.IsUsernameExists(l.svcCtx, in.Username); err != nil {
		return &identity.RegisterResp{}, err
	}

	// create user
	pwdEncrypted, err := tool.EncipherPassword(in.Password)
	if err != nil {
		logx.Errorf("failed to encipher password: %v", err)
		return &identity.RegisterResp{}, errorx.ErrSystemError
	}

	var User = models.User{
		Username: in.Username,
		Password: pwdEncrypted,
		Nickname: userutil.GenerateNickname(),
	}

	if res := l.svcCtx.Db.Create(&User); res.Error != nil || res.RowsAffected == 0 {
		logx.Errorf("failed to create user: %v", err)
		return &identity.RegisterResp{}, errorx.ErrSystemError
	}

	// create token
	token, err := userutil.GenerateJwt(l.svcCtx, User.ID, in.Username)
	if err != nil {
		logx.Errorf("failed to generate jwt: %v", err)
		return &identity.RegisterResp{}, errorx.ErrSystemError
	}

	return &identity.RegisterResp{
		Response: &identity.Resp{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: User.ID,
		Token:  token,
	}, nil
}
