package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
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
	PreMsgTime := in.PreMsgTime

	var Messages []models.Chat
	if res := l.svcCtx.Db.Model(&models.Chat{}).
		Where("from_user_id = ? AND to_user_id = ?", in.FromUserId, in.ToUserId).
		Where("create_at < ?", PreMsgTime).
		Find(&Messages); errors.Is(res.Error, gorm.ErrRecordNotFound) {

		return &chat.ChatMessageResp{}, nil
	} else if res.Error != nil {
		logx.Errorf("get chat message failed: %v", res.Error)

		// to/do bug：只是查询了单方的信息，应该查询双方信息
		return &chat.ChatMessageResp{}, res.Error
	}
	// fix

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
