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
		return &types.FriendListResp{}, err
	}

	var userList []types.FriendUser
	// get user detail information
	for _, targetUserId := range friendRpc.UserIdList {
		// get user information
		userInfo, err := rpcutil.GetUserInfo(l.svcCtx, l.ctx, req.UserId, targetUserId)
		if err != nil {
			logx.Errorf("get user information failed: %v", err)

			return &types.FriendListResp{}, nil
		}

		// get last chat message
		chatLastMessageRpc, err := l.svcCtx.ChatRpc.ChatLastMessage(l.ctx, &chat.ChatLastMessageReq{
			ToUserId:   req.UserId,
			FromUserId: targetUserId,
		})
		if err != nil {
			logx.Errorf("get last chat message failed: %v", err)

			return &types.FriendListResp{}, err
		}

		if chatLastMessageRpc.Message == nil {
			userList = append(userList, types.FriendUser{
				Message: "来打个招呼吧",
				MsgType: int64(1),
				User:    userInfo,
			})
			continue
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
