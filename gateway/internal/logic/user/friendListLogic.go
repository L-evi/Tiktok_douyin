package user

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/service/chat/types/chat"
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

	if req.UserId != l.ctx.Value("user_id").(int64) {
		return &types.FriendListResp{}, nil
	}

	// get friendIdList from friend rpc
	var friendRpc *user.FriendListResp
	if friendRpc, err = l.svcCtx.UserRpc.FriendList(l.ctx, &user.FriendListReq{
		UserId: req.UserId,
	}); err != nil {

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

		// get last chat message
		chatLastMessageRpc, err := l.svcCtx.ChatRpc.ChatLastMessage(l.ctx, &chat.ChatLastMessageReq{
			ToUserId:   req.UserId,
			FromUserId: value,
		})
		if err != nil {
			logx.Errorf("get last chat message failed: %v", err)

			return &types.FriendListResp{}, err
		}

		// 判断信息类型 0: send message, 1: receive message
		msgType := int64(0)
		if chatLastMessageRpc.Message.FromUserId == req.UserId {
			msgType = int64(1)
		}

		userList = append(userList, types.FriendUser{
			Message: chatLastMessageRpc.Message.Content,
			MsgType: msgType,
			User:    userInfo,
		})
	}

	return &types.FriendListResp{
		Resp:     errx.SUCCESS_RESP,
		UserList: userList,
	}, nil
}
