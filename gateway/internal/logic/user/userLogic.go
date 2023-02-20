package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
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

	userInfo, err := rpcutil.GetUserInfo(l.svcCtx, l.ctx, userId, req.UserId)
	if err != nil {
		logx.Errorf("get user information failed: %v", err)

		return &types.UserResp{}, errorx.ErrDatabaseError
	}
	return &types.UserResp{
		Resp: errx.SUCCESS_RESP,
		User: userInfo,
	}, nil
}
