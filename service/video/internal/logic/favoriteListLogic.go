package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"train-tiktok/common/errorx"
	tool2 "train-tiktok/common/tool"
	"train-tiktok/service/video/common/rediskeyutil"
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
	// check user exists
	if exists, err := tool2.CheckUserExist(l.ctx, l.svcCtx.IdentityRpc, in.UserId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query user: %v", err)

		return nil, errorx.ErrSystemError
	} else if !exists {
		return nil, errorx.ErrUserNotFound
	}

	rdb := l.svcCtx.Rdb

	_userKey := rediskeyutil.NewKeys(l.svcCtx.Config.RedisConf.Prefix).GetUserKey(in.UserId)

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
			continue
		}

		var myVideo models.Video

		if err := l.svcCtx.Db.Where(&models.Video{ID: videoId}).First(&myVideo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("failed to get video, err: %v", err)
			rdb.ZRem(l.ctx, _userKey, v)
			continue
		} else if err != nil {
			logx.Errorf("failed to get video, err: %v", err)

			return &video.FavoriteListResp{}, errorx.ErrDatabaseError
		}

		myVideo.PlayUrl, myVideo.CoverUrl = tool.HandleVideoUrl(l.svcCtx, myVideo.Position, myVideo.PlayUrl, myVideo.CoverUrl)
		// TO/DO
		favoriteList = append(favoriteList, &video.VideoX{
			Id:       myVideo.ID,
			UserId:   myVideo.UserID,
			PlayUrl:  myVideo.PlayUrl,
			CoverUrl: myVideo.CoverUrl,
			Title:    myVideo.Title,
		})
	}
	// consult success
	logx.WithContext(l.ctx).Info(favoriteList)
	return &video.FavoriteListResp{
		VideoList: favoriteList,
	}, nil
}
