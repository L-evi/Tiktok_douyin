package user

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/user/types/user"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionReq) (resp *types.RelationActionResp, err error) {
	_userId := l.ctx.Value("user_id").(int64)

	if _, err := l.svcCtx.UserRpc.RelationAct(l.ctx, &user.RelationActReq{
		UserId:   _userId,
		TargetId: req.ToUserId,
		Action:   req.ActionTyp,
	}); err != nil {
		return &types.RelationActionResp{}, err
	}

	return &types.RelationActionResp{
		Resp: errx.SUCCESS_RESP,
	}, nil
}
