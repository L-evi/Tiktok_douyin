package user

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	"train-tiktok/service/user/types/user"

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

func (l *UserLogic) User(req *types.UserReq) (resp *types.UserResp, err error) {
	var userId int64
	var isLogin bool
	if isLogin = l.ctx.Value("is_login").(bool); isLogin {
		userId = l.ctx.Value("user_id").(int64)
	}

	rpcResp, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
		UserId:   userId,
		TargetId: req.UserId,
	})
	if err != nil {
		return &types.UserResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	return &types.UserResp{
		Resp: errx.SUCCESS_RESP,
		User: types.User{
			Id:            req.UserId,
			Name:          rpcResp.Name,
			FollowCount:   *rpcResp.FollowCount,
			FollowerCount: *rpcResp.FollowerCount,
			IsFollow:      rpcResp.IsFollow,
		},
	}, nil
}
