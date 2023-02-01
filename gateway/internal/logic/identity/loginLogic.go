package identity

import (
	"context"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/identity/types/identity"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *identityclient.LoginResp, err error) {
	return l.svcCtx.IdentityRpc.Login(l.ctx, &identity.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
}
