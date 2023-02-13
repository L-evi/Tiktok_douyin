package logic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"train-tiktok/service/video/common/tool"
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
	rdb := l.svcCtx.Rdb

	_userKey := fmt.Sprintf("%s:favorite_user:%d", l.svcCtx.Config.RedisConf.Prefix, in.UserId)

	var err error
	var userFavoriteList []string
	if userFavoriteList, err = rdb.ZRange(l.ctx, _userKey, 0, -1).Result(); err != nil {
		logx.Errorf("failed to get user favorite list, err: %v", err)

		return &video.FavoriteListResp{}, err
	}

	var favoriteList []*video.VideoX
	favoriteList = make([]*video.VideoX, 0, len(userFavoriteList))
	for _, v := range userFavoriteList {
		// get favorite video information
		var videoId int64
		var err error
		if videoId, err = strconv.ParseInt(v, 10, 64); err != nil {
			logx.Errorf("failed to convert video id, err: %v", err)
			rdb.ZRem(l.ctx, _userKey, v)

			return &video.FavoriteListResp{}, err
		}
		var myVideo models.Video
		result := l.svcCtx.Db.Where(&models.Video{ID: videoId}).First(&myVideo)
		if result.Error != nil {
			logx.Errorf("failed to get video, err: %v", result.Error)

			return &video.FavoriteListResp{}, result.Error
		}

		// TODO
		favoriteList = append(favoriteList, &video.VideoX{
			Id:       myVideo.ID,
			UserId:   myVideo.UserID,
			PlayUrl:  tool.GetFullPlayUrl(l.svcCtx, myVideo.Position, myVideo.PlayUrl),
			CoverUrl: tool.GetFullCoverUrl(l.svcCtx, myVideo.Position, myVideo.CoverUrl),
			Title:    myVideo.Title,
		})
	}
	// consult success
	log.Println(favoriteList)
	return &video.FavoriteListResp{
		VideoList: favoriteList,
	}, nil
}
