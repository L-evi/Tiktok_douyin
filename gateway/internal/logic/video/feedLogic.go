package video

import (
	"context"
	"train-tiktok/common/errorx"
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
		// 点赞
		var isFavor = false
		var CommentCount int64
		var FavoriteCount int64

		if isLogin {
			if favorite, err := isFavorite(l.svcCtx, l.ctx, userId, v.Id); err != nil {
				return &types.FeedResp{}, errorx.ErrSystemError
			} else {
				isFavor = favorite
			}
		}
		if favorCount, err := getFavoriteCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FeedResp{}, errorx.ErrSystemError
		} else {
			FavoriteCount = favorCount
		}
		if _commentCount, err := getCommentCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FeedResp{}, errorx.ErrSystemError
		} else {
			CommentCount = _commentCount
		}

		videoList = append(videoList, types.Video{
			Id:            v.Id,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    isFavor,
		})
	}
	logx.WithContext(l.ctx).Infof("videoList: %v", videoList)

	return &types.FeedResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}

func isFavorite(c *svc.ServiceContext, ctx context.Context, userId int64, videoId int64) (bool, error) {
	var err error
	var resp *video.IsFavoriteResp
	if resp, err = c.VideoRpc.IsFavorite(ctx, &video.IsFavoriteReq{
		UserId:  userId,
		VideoId: videoId,
	}); err != nil {
		return false, err
	}
	return resp.IsFavorite, nil
}

func getFavoriteCount(c *svc.ServiceContext, ctx context.Context, videoId int64) (int64, error) {
	var err error
	var resp *video.FavoriteCountResp
	if resp, err = c.VideoRpc.FavoriteCount(ctx, &video.FavoriteCountReq{
		VideoId: videoId,
	}); err != nil {
		return 0, err
	}
	return resp.FavoriteCount, nil
}

func getCommentCount(c *svc.ServiceContext, ctx context.Context, videoId int64) (int64, error) {
	var err error
	var resp *video.CommentCountResp
	if resp, err = c.VideoRpc.CommentCount(ctx, &video.CommentCountReq{
		VideoId: videoId,
	}); err != nil {
		return 0, err
	}
	return resp.CommentCount, nil
}
