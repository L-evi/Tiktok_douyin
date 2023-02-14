package tool

import (
	"fmt"
	"gorm.io/gorm"
	"train-tiktok/service/video/internal/svc"
	"train-tiktok/service/video/models"
)

func HandleVideoUrl(svcCtx *svc.ServiceContext, position, playUrl, coverUrl string) (string, string) {
	switch position {
	case "local":
		playUrl = fmt.Sprintf("%s/%s", svcCtx.StorageBaseUrl.Local, playUrl)
		coverUrl = fmt.Sprintf("%s/%s", svcCtx.StorageBaseUrl.Local, coverUrl)
		break
	default:
		break
	}
	return playUrl, coverUrl
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
