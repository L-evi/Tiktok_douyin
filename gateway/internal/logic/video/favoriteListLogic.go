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

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {

	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListReq) (resp *types.FovoriteListResp, err error) {
	// sent to rpc to consult
	var userId int64
	var isLogin bool
	if isLogin = l.ctx.Value("is_login").(bool); isLogin {
		userId = l.ctx.Value("user_id").(int64)
	}

	var rpcResp *video.FavoriteListResp
	if rpcResp, err = l.svcCtx.VideoRpc.FavoriteList(l.ctx, &video.FavoriteListReq{
		UserId: req.UserId,
	}); err != nil {
		logx.Errorf("get favorite list failed, err: %v", err)

		return &types.FovoriteListResp{
			Resp:      errx.HandleRpcErr(err),
			VideoList: nil,
		}, nil
	}

	// consult success
	var videoList []types.Video
	for _, v := range rpcResp.VideoList {
		// 点赞
		var isFavor = false
		var CommentCount int64
		var FavoriteCount int64

		if isLogin {
			if isFavor, err = rpcutil.IsFavorite(l.svcCtx, l.ctx, userId, v.Id); err != nil {
				return &types.FovoriteListResp{}, errorx.ErrSystemError
			}
		}
		if FavoriteCount, err = rpcutil.GetFavoriteCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FovoriteListResp{}, errorx.ErrSystemError
		}
		if CommentCount, err = rpcutil.GetCommentCount(l.svcCtx, l.ctx, v.Id); err != nil {
			return &types.FovoriteListResp{}, errorx.ErrSystemError
		}
		// getUserInfo
		var userInfo types.User
		if !isLogin {
			userId = v.UserId
		}
		if userInfo, err = rpcutil.GetUserInfo(l.svcCtx, l.ctx, userId, v.UserId); err != nil {
			return &types.FovoriteListResp{}, errorx.ErrSystemError
		}

		videoList = append(videoList, types.Video{
			Id: v.Id,
			Author: types.User{
				Id:            userId,
				Name:          userInfo.Name,
				FollowCount:   userInfo.FollowCount,
				FollowerCount: userInfo.FollowerCount,
				IsFollow:      userInfo.IsFollow,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    isFavor,
			Title:         v.Title,
		})
	}

	return &types.FovoriteListResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}
