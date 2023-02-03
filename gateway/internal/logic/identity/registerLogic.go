package identity

import (
	"context"
	"train-tiktok/gateway/common/errx"
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

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	rpcResp, err := l.svcCtx.IdentityRpc.Register(l.ctx, &identity.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return &types.RegisterResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	return &types.RegisterResp{
		Resp:   errx.SUCCESS_RESP,
		UserId: rpcResp.UserId,
		Token:  rpcResp.Token,
	}, nil
}
