package logic

import (
	"context"

	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLogic) User(in *user.UserReq) (*user.UserResp, error) {
	return &user.UserResp{
		Name: in.Username,
	}, nil
}
