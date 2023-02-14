package video

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
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
	rpcResp, err := l.svcCtx.VideoRpc.Feed(l.ctx, &video.FeedReq{
		LatestTime: req.LatestTime,
	})
	if err != nil {
		return &types.FeedResp{}, err
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
			if isFavor, err = rpcutil.IsFavorite(l.svcCtx, l.ctx, userId, v.Id); err != nil {
				return &types.FeedResp{}, errorx.ErrSystemError
			}
		}
		if FavoriteCount, err = rpcutil.GetFavoriteCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FeedResp{}, errorx.ErrSystemError
		}
		if CommentCount, err = rpcutil.GetCommentCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FeedResp{}, errorx.ErrSystemError
		}

		// getUserInfo
		var userInfo types.User
		if !isLogin {
			userId = 0 // isFollow 将返回 false
		}
		if userInfo, err = rpcutil.GetUserInfo(l.svcCtx, l.ctx, userId, v.UserId); err != nil {
			return &types.FeedResp{}, errorx.ErrSystemError
		}

		videoList = append(videoList, types.Video{
			Id:            v.Id,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    isFavor,
			Author:        userInfo,
		})
	}
	logx.WithContext(l.ctx).Infof("videoList: %v", videoList)

	return &types.FeedResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}
