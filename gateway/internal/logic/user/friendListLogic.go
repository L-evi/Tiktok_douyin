package user

import (
	"context"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/service/user/user"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// get friendIdList from friend rpc
	friendRpc, err := l.svcCtx.UserRpc.FriendList(l.ctx, &user.FriendListReq{
		UserId: req.UserId,
	})
	if err != nil {

		return &types.FriendListResp{}, nil
	}
	friendIdList := friendRpc.UserIdList

	var userList []types.FriendUser
	// get user detail information
	for _, value := range friendIdList {
		// get user information
		userInfo, err := rpcutil.GetUserInfo(l.svcCtx, l.ctx, req.UserId, value)
		if err != nil {
			logx.Errorf("get user information failed: %v", err)

			return &types.FriendListResp{}, nil
		}

		friendUser := types.FriendUser{
			// todo: 在聊天中添加一个查询最近的消息的请求
			Message: "",
			MsgType: 0,
			User:    userInfo,
		}
		userList = append(userList, friendUser)
	}

	return &types.FriendListResp{
		Resp:     types.Resp{},
		UserList: userList,
	}, nil
}
