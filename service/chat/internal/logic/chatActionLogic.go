package logic

import (
	"context"
	"time"
	"train-tiktok/service/chat/internal/svc"
	"train-tiktok/service/chat/models"
	"train-tiktok/service/chat/types/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatActionLogic {
	return &ChatActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatActionLogic) ChatAction(in *chat.ChatActionReq) (*chat.CharActionResp, error) {
	// in.ActionType = 1 sent message
	if in.ActionType == 1 {
		var chatMessage = models.Chat{
			FromUserId: in.FromUserId,
			ToUserId:   in.ToUserId,
			CreateAt:   time.Now().Unix(),
			Content:    in.Content,
		}
		// create chat into database
		if res := l.svcCtx.Db.Create(&chatMessage); res.Error != nil {
			logx.WithContext(l.ctx).Errorf("failed to create chat: %v", res.Error)

			return &chat.CharActionResp{}, res.Error
		}
	}

	return &chat.CharActionResp{}, nil
}
