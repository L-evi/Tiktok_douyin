package tool

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"train-tiktok/common/position"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/models"
)

// HandleVideoUrl 处理视频的url
// position: local, cos
// if start with http, return directly
func HandleVideoUrl(svcCtx *svc.ServiceContext, posit, playUrl, coverUrl string) (string, string) {

	switch posit {
	case position.LOCAL:
		playUrl = parseVideoUrl(svcCtx.StorageBaseUrl.Local, playUrl)
		coverUrl = parseVideoUrl(svcCtx.StorageBaseUrl.Local, coverUrl)
		break
	case position.COS:
		playUrl = parseVideoUrl(svcCtx.StorageBaseUrl.Cos, playUrl)
		coverUrl = parseVideoUrl(svcCtx.StorageBaseUrl.Cos, coverUrl)
	default:
		break
	}
	return playUrl, coverUrl
}

// 合成视频的url. 如果是http开头, 则直接返回, 否则拼接base
func parseVideoUrl(base, url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	if strings.HasPrefix(url, "/") || strings.HasPrefix(url, "./") {
		url = strings.TrimLeft(url, "/")
		url = strings.TrimLeft(url, "./")
	}
	return fmt.Sprintf("%s/%s", base, url)
}

func CheckVideoExists(db *gorm.DB, videoId int64) (bool, error) {
	var count int64
	if err := db.Model(&models.Video{}).
		Where(&models.Video{ID: videoId}).Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetVideoUserId(db *gorm.DB, videoId int64) (int64, error) {
	var video models.Video
	if err := db.Model(&models.Video{}).
		Where(&models.Video{ID: videoId}).First(&video).
		Error; err != nil {

		return 0, err
	}

	return video.UserID, nil
}
