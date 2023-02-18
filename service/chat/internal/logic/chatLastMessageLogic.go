package logic

import (
	"context"
	"train-tiktok/service/chat/internal/svc"
	"train-tiktok/service/chat/models"
	"train-tiktok/service/chat/types/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLastMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatLastMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLastMessageLogic {
	return &ChatLastMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatLastMessageLogic) ChatLastMessage(in *chat.ChatLastMessageReq) (*chat.ChatLastMessageResp, error) {
	var message *chat.Message
	if err := l.svcCtx.Db.
		Limit(1).
		Order("CreateAt desc").
		Where(&models.Chat{FromUserId: in.FromUserId, ToUserId: in.ToUserId}).
		Or(&models.Chat{FromUserId: in.ToUserId, ToUserId: in.FromUserId}).
		Find(message).Error; err != nil {
		logx.Errorf("get last chat message failed: %v", err)

		return &chat.ChatLastMessageResp{}, err
	}

	return &chat.ChatLastMessageResp{
		Resp:    nil,
		Message: message,
	}, nil
}
