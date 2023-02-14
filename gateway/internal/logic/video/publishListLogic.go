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

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	userId := l.ctx.Value("user_id").(int64)

	var rpcResp *video.PublishListResp
	if rpcResp, err = l.svcCtx.VideoRpc.PublishList(l.ctx, &video.PublishListReq{
		UserId: req.UserId, // 此处传入的是 query 中的 userId, 也就是被查看的用户的 id
	}); err != nil {
		return &types.PublishListResp{}, err
	}

	var videoList []types.Video
	videoList = make([]types.Video, 0, len(rpcResp.VideoList))
	for _, v := range rpcResp.VideoList {
		// 点赞
		var favorite = false
		var commentCount int64
		var favorCount int64

		if favorite, err = rpcutil.IsFavorite(l.svcCtx, l.ctx, userId, v.Id); err != nil {
			return &types.PublishListResp{}, errorx.ErrSystemError
		}
		if favorCount, err = rpcutil.GetFavoriteCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.PublishListResp{}, errorx.ErrSystemError
		}
		if commentCount, err = rpcutil.GetCommentCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.PublishListResp{}, errorx.ErrSystemError
		}

		// getUserInfo
		var userInfo types.User
		if userInfo, err = rpcutil.GetUserInfo(l.svcCtx, l.ctx, userId, v.UserId); err != nil {
			return &types.PublishListResp{}, errorx.ErrSystemError
		}

		videoList = append(videoList, types.Video{
			Id:            v.Id,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: favorCount,
			CommentCount:  commentCount,
			IsFavorite:    favorite,
			Author:        userInfo,
		})
	}
	logx.WithContext(l.ctx).Infof("publishlist: %v", videoList)

	return &types.PublishListResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}
