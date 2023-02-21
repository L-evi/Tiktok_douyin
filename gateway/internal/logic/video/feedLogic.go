package video

import (
	"context"
	"time"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"
	videoErrx "train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/types/video"

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
	var rpcResp *video.FeedResp

	if rpcResp, err = l.svcCtx.VideoRpc.Feed(l.ctx, &video.FeedReq{
		LatestTime: req.LatestTime,
	}); errorx.IsRpcError(err, videoErrx.ErrNoLatestVideo) || len(rpcResp.VideoList) == 0 {
		// 没有更新的视频了
		// loop
		//return NewFeedLogic(l.ctx, l.svcCtx).Feed(&types.FeedReq{
		//	LatestTime: time.Now().Unix(),
		//})
		return &types.FeedResp{
			NextTime: time.Now().Unix(),
		}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("Feed rpc error: %v", err)

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
	logx.WithContext(l.ctx).Debugf("videoList: %v", videoList)
	logx.WithContext(l.ctx).Debugf("nextTime: %v", *rpcResp.NextTime)
	return &types.FeedResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
		NextTime:  *rpcResp.NextTime,
	}, nil
}
