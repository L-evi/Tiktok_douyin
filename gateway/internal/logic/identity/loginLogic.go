package identity

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	"train-tiktok/service/identity/types/identity"

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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.IdentityRpc.Login(l.ctx, &identity.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return &types.LoginResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	return &types.LoginResp{
		Resp:   errx.SUCCESS_RESP,
		UserId: rpcResp.UserId,
		Token:  rpcResp.Token,
	}, nil
}
