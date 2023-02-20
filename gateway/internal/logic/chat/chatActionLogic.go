package chat

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/common/tool"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/chat/types/chat"
	"train-tiktok/service/user/types/user"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatActionLogic {
	return &ChatActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatActionLogic) ChatAction(req *types.ChatActionReq) (resp *types.ChatActionResp, err error) {
	fromUserId := l.ctx.Value("user_id").(int64)

	if exists, err := tool.CheckUserExist(l.ctx, l.svcCtx.IdentityRpc, fromUserId); err != nil {
		return &types.ChatActionResp{}, errorx.ErrSystemError
	} else if !exists {
		return &types.ChatActionResp{}, errorx.ErrUserNotFound
	}

	// check is friend
	if rpcResp, err := l.svcCtx.UserRpc.IsFriend(l.ctx, &user.IsFriendReq{
		UserId:   fromUserId,
		TargetId: req.ToUserId,
	}); err != nil {
		return &types.ChatActionResp{}, err
	} else if !rpcResp.IsFriend {
		logx.Errorf("not friend: %s <=> %s -> %v", fromUserId, req.ToUserId, err)
		return &types.ChatActionResp{}, errorx.ErrNotFriend
	}

	if _, err = l.svcCtx.ChatRpc.ChatAction(l.ctx, &chat.ChatActionReq{
		FromUserId: fromUserId,
		ToUserId:   req.ToUserId,
		Content:    req.Content,
		ActionType: req.ActionType,
	}); err != nil {

		return &types.ChatActionResp{}, err
	}

	return &types.ChatActionResp{
		Resp: errx.SUCCESS_RESP,
	}, nil
}
