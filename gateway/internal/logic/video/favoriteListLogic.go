package video

import (
	"context"
	"log"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/user/types/user"
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
	rpcResp, err := l.svcCtx.VideoRpc.FavoriteList(l.ctx, &video.FavoriteListReq{
		UserId: req.UserId,
	})
	// consult failed
	if err != nil {
		log.Printf("get favorite list failed, err: %v", err)
		return &types.FovoriteListResp{
			Resp:      errx.HandleRpcErr(err),
			VideoList: nil,
		}, nil
	}
	// consult success
	var videoList []types.Video
	for _, v := range rpcResp.VideoList {
		var userId = v.UserId
		// get author info from user-rpc-service
		userRpcResp, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
			UserId:   l.ctx.Value("user_id").(int64),
			TargetId: userId,
		})
		// consult failed
		if err != nil {
			log.Printf("get user information failed: %v", err)
			return &types.FovoriteListResp{
				Resp:      errx.HandleRpcErr(err),
				VideoList: nil,
			}, nil
		}
		// consult success
		var author = types.User{
			Id:            userId,
			Name:          userRpcResp.Name,
			FollowCount:   *userRpcResp.FollowCount,
			FollowerCount: *userRpcResp.FollowerCount,
			IsFollow:      userRpcResp.IsFollow,
		}
		videoList = append(videoList, types.Video{
			Id:            v.Id,
			Author:        author,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsForavite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	return &types.FovoriteListResp{
		Resp:      errx.SUCCESS_RESP,
		VideoList: videoList,
	}, nil
}
