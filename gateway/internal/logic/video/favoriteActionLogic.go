package video

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/video/types/video"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionReq) (resp *types.FavoriteActionResp, err error) {
	// sent to rpc to consult
	_, err = l.svcCtx.VideoRpc.FavoriteAction(l.ctx, &video.FavoriteActionReq{
		UserId:     l.ctx.Value("user_id").(int64),
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	// consult failed
	if err != nil {
		return &types.FavoriteActionResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}
	// consult success
	return &types.FavoriteActionResp{
		Resp: errx.SUCCESS_RESP,
	}, nil
}
