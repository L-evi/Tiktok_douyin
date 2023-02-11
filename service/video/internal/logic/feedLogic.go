package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
	"train-tiktok/service/video/common/errx"
	"train-tiktok/service/video/models"

	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {

	return &FeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeedLogic) Feed(in *video.FeedReq) (*video.FeedResp, error) {
	lastTime := in.LatestTime // 可能为空 > 传入当前时间戳

	// lastTime = 0 时，返回最新的10条数据
	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}

	// query
	var videos []models.Video
	if err := l.svcCtx.Db.Model(&models.Video{}).Where("create_at < ?", lastTime).
		Limit(10).Order("create_at desc").
		Find(&videos).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		return &video.FeedResp{}, errx.ErrNoLatestVideo
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("FeedLogic Feed sql err: %v", err)

		return &video.FeedResp{}, err
	}

	// 处理数据
	// 获取 videos 内的最老 时间 并生成 videoList
	var videoList []*video.Video
	var nextTime = lastTime
	videoList = make([]*video.Video, 0, len(videos))

	for _, v := range videos {
		if v.CreateAt < nextTime {
			nextTime = v.CreateAt
		}

		position := v.Position // 视频 存储节点  cos or  local
		switch position {
		case "local":
			v.PlayUrl = l.svcCtx.StorageBaseUrl.Local + v.PlayUrl
			break
		default:
			break
		}

		// insert videoList
		videoList = append(videoList, &video.Video{
			Id:       v.ID,
			UserId:   v.UserID,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		})
	}

	return &video.FeedResp{
		NextTime:  &nextTime,
		VideoList: videoList,
	}, nil
}
