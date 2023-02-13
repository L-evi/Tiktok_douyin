package logic

import (
	"context"
	"train-tiktok/service/chat/models"

	"train-tiktok/service/chat/internal/svc"
	"train-tiktok/service/chat/types/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageLogic {
	return &ChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatMessageLogic) ChatMessage(in *chat.ChatMessageReq) (*chat.ChatMessageResp, error) {
	var Messages []models.Chat
	if res := l.svcCtx.Db.Model(&models.Chat{}).Where(&models.Chat{FromUserId: in.FromUserId, ToUserId: in.ToUserId}).Find(&Messages); res.Error != nil {
		logx.Errorf("get chat message failed: %v", res.Error)

		return &chat.ChatMessageResp{}, res.Error
	}
	var list []*chat.Message
	for _, v := range Messages {
		list = append(list, &chat.Message{
			Id:         v.ID,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    v.Content,
			CreateTime: v.CreateAt,
		})
	}

	return &chat.ChatMessageResp{
		MessageList: list,
	}, nil
}
