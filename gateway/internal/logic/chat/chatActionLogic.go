package chat

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/chat/types/chat"

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
	_, err = l.svcCtx.ChatRpc.ChatAction(l.ctx, &chat.ChatActionReq{
		FromUserId: l.ctx.Value("user_id").(int64),
		ToUserId:   req.ToUserId,
		Content:    req.Content,
		ActionType: req.ActionType,
	})
	if err != nil {

		return &types.ChatActionResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	return &types.ChatActionResp{
		Resp: errx.SUCCESS_RESP,
	}, nil
}
