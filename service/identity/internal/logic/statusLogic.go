package logic

import (
	"context"
	"train-tiktok/service/identity/common/errx"
	"train-tiktok/service/identity/common/userutil"
	"train-tiktok/service/identity/models"

	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/types/identity"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusLogic {
	return &StatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Status check user login status
func (l *StatusLogic) Status(in *identity.StatusReq) (*identity.StatusResp, error) {
	if in.Token == "" {
		return &identity.StatusResp{
			IsLogin: false,
		}, errx.ErrTokenInvalid
	}

	var err error
	var _user models.User
	if _user, err = userutil.CheckPermission(l.svcCtx, in.Token); err != nil {
		return &identity.StatusResp{}, err
	}

	return &identity.StatusResp{
		IsLogin:  true,
		UserId:   _user.ID,
		Username: _user.Username,
	}, nil
}
