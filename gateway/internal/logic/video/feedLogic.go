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

	// how to know if user_id in context
	var userId int64
	var isLogin bool
	if isLogin = l.ctx.Value("is_login").(bool); isLogin {
		userId = l.ctx.Value("user_id").(int64)
	}

	var videoList []types.Video
	videoList = make([]types.Video, 0, len(rpcResp.VideoList))
	for _, v := range rpcResp.VideoList {
		// TODO 点赞
		isFavo := false
		if isLogin {
			favo, err := isFavorite(l.svcCtx, userId, v.Id)
			// if
		}
		videoList = append(videoList, types.Video{
			Id:            v.Id,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		})
	}

	return &types.FeedResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}

func isFavorite(c *svc.ServiceContext, userId int64, videoId int64) (bool, error) {
	// c.UserRpc
	return false, nil
}
