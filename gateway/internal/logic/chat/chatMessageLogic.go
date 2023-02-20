package chat

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/chat/types/chat"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageLogic {
	return &ChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatMessageLogic) ChatMessage(req *types.ChatMessageReq) (resp *types.ChatMessageResp, err error) {

	var rpcResp *chat.ChatMessageResp
	if rpcResp, err = l.svcCtx.ChatRpc.ChatMessage(l.ctx, &chat.ChatMessageReq{
		FromUserId: l.ctx.Value("user_id").(int64),
		ToUserId:   req.ToUserId,
		PreMsgTime: req.PreMsgTime,
	}); err != nil {

		return &types.ChatMessageResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	var Messages []types.Message
	for _, v := range rpcResp.MessageList {
		Messages = append(Messages, types.Message{
			Id:         v.Id,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    v.Content,
			CreateTime: v.CreateTime / 1000,
		})
	}

	return &types.ChatMessageResp{
		Resp:        errx.SUCCESS_RESP,
		MessageList: Messages,
	}, nil
}
