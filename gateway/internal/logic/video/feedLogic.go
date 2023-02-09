package video

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/video/types/video"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	rpcResp, err := l.svcCtx.VideoRpc.Feed(l.ctx, &video.FeedReq{})
	if err != nil {
		return &types.FeedResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	videoList :=

	return &types.FeedResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: rpcResp.VideoList,
	}, nil
}
