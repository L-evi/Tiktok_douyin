package logic

import (
	"context"

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
	return &identity.RegisterResp{
		Response: &identity.Resp{
			StatusCode: 200,
			StatusMsg:  "sss",
		},
		UserId: 0,
		Token:  in.Username,
	}, nil
}
