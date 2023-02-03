package user

import (
	"context"
	"train-tiktok/service/user/types/user"
	"train-tiktok/service/user/userclient"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.UserReq) (resp *userclient.UserResp, err error) {
	return l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
		UserId:   l.ctx.Value("userId").(int64),
		TargetId: req.User_id,
	})
}
