package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
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
	var message models.Chat
	if err := l.svcCtx.Db.Model(&models.Chat{}).
		Where(&models.Chat{FromUserId: in.FromUserId, ToUserId: in.ToUserId}).
		Or(&models.Chat{FromUserId: in.ToUserId, ToUserId: in.FromUserId}).
		First(&message).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		return &chat.ChatLastMessageResp{}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("get last chat message failed: %s", err)

		return &chat.ChatLastMessageResp{}, errorx.ErrDatabaseError
	}

	return &chat.ChatLastMessageResp{
		Message: &chat.Message{
			Id:         message.ID,
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreateAt,
		},
	}, nil
}
