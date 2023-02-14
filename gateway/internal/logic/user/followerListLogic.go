package user

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/service/user/types/user"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowerListLogic) FollowerList(req *types.FansListReq) (resp *types.FansListResp, err error) {
	var userId int64
	var isLogin bool
	if isLogin = l.ctx.Value("is_login").(bool); isLogin {
		userId = l.ctx.Value("user_id").(int64)
	}

	var rpcResp *user.FollowerListResp
	if rpcResp, err = l.svcCtx.UserRpc.FollowerList(l.ctx, &user.FollowerListReq{
		UserId: req.UserId,
	}); err != nil {
		return &types.FansListResp{}, err
	}

	var users []types.User
	users = make([]types.User, 0, len(rpcResp.UserIds))
	for _, rpcUserId := range rpcResp.UserIds {
		// getUserInfo
		var userInfo types.User
		if !isLogin {
			userId = 0 // isFollow 将返回 false
		}
		if userInfo, err = rpcutil.GetUserInfo(l.svcCtx, l.ctx, userId, rpcUserId); err != nil {
			return &types.FansListResp{}, errorx.ErrSystemError
		}

		users = append(users, types.User{
			Id:            rpcUserId,
			Name:          userInfo.Name,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      userInfo.IsFollow,
		})
	}

	return &types.FansListResp{
		Resp:     errx.SUCCESS_RESP,
		UserList: users,
	}, nil
}
