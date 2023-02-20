package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
	"train-tiktok/service/video/common/tool"
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

		// 没有更新的数据 //  errx.ErrNoLatestVideo
		return &video.FeedResp{}, nil
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("FeedLogic Feed sql err: %v", err)

		return &video.FeedResp{}, err
	}

	logx.WithContext(l.ctx).Debugf("FeedLogic Feed videos: %v", videos)
	if len(videos) == 0 {
		// 没有更新的数据 //  errx.ErrNoLatestVideo
		return &video.FeedResp{}, nil
	}

	// 处理数据
	// 获取 videos 内的最老 时间 并生成 videoList
	var videoList []*video.VideoX
	var nextTime = lastTime
	videoList = make([]*video.VideoX, 0, len(videos))

	for _, v := range videos {
		if v.CreateAt < nextTime {
			nextTime = v.CreateAt
		}

		v.PlayUrl, v.CoverUrl = tool.HandleVideoUrl(l.svcCtx, v.Position, v.PlayUrl, v.CoverUrl)

		// insert videoList
		videoList = append(videoList, &video.VideoX{
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
