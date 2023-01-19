package identity

import (
	"context"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/identity/types/identity"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *identityclient.RegisterResp, err error) {
	return l.svcCtx.IdentityRpc.Register(l.ctx, &identity.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})
}