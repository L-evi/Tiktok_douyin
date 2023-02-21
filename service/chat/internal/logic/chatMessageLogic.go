package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
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
	if err := l.svcCtx.Db.Model(&models.Chat{}).Select([]string{"id", "from_user_id", "to_user_id", "content", "create_at"}).
		Where("(from_user_id = ? and to_user_id = ? or from_user_id = ? and to_user_id = ?)", in.FromUserId, in.ToUserId, in.ToUserId, in.FromUserId).
		Where("create_at > ? ", in.PreMsgTime).
		Order("create_at desc").Limit(30).
		Find(&Messages).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		return &chat.ChatMessageResp{}, nil
	} else if err != nil {
		logx.Errorf("get chat message failed: %v", err)

		// to/do bug：只是查询了单方的信息，应该查询双方信息
		return &chat.ChatMessageResp{}, errorx.ErrDatabaseError
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
