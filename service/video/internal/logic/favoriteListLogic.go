package logic

import (
	"context"
	UserModels "train-tiktok/service/user/models"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteListLogic) FavoriteList(in *video.FavoriteListReq) (*video.FavoriteListResp, error) {
	// get user_favorite list
	var userFavoriteList []UserModels.UserFavorite
	res := l.svcCtx.Db.Where(&UserModels.UserFavorite{UserId: in.UserId}).Find(&userFavoriteList)
	if res.Error != nil {
		logx.Errorf("failed to get user favorite, err: %v", res.Error)
		return &video.FavoriteListResp{}, res.Error
	}

	var favoriteList []*video.FavoriteVideo
	favoriteList = make([]*video.FavoriteVideo, len(userFavoriteList))

	for _, v := range userFavoriteList {
		// get favorite video information
		var myVideo models.Video
		result := l.svcCtx.Db.Where(&models.Video{ID: v.VideoId}).First(&myVideo)
		if result.Error != nil {
			logx.Errorf("failed to get video, err: %v", result.Error)
			return &video.FavoriteListResp{}, result.Error
		}

		favoriteList = append(favoriteList, &video.FavoriteVideo{
			Id:            myVideo.ID,
			UserId:        myVideo.UserID,
			PlayUrl:       myVideo.PlayUrl,
			CoverUrl:      myVideo.CoverUrl,
			IsFavorite:    true,
			CommentCount:  1,
			FavoriteCount: 1,
			Title:         myVideo.Title,
		})
	}
	// consult success
	return &video.FavoriteListResp{
		VideoList: favoriteList,
	}, nil
}
