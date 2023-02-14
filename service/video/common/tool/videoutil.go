package tool

import (
	"fmt"
	"strings"
	"train-tiktok/service/video/internal/svc"
)

func GetFullPlayUrl(svcCtx *svc.ServiceContext, position, playUrl string) string {
	if strings.HasPrefix(playUrl, "http") {
		return playUrl
	}
	switch position {
	case "local":
		return fmt.Sprintf("%s/%s", svcCtx.StorageBaseUrl.Local, playUrl)
	default:
		return playUrl
	}
}

func GetFullCoverUrl(svcCtx *svc.ServiceContext, position, coverUrl string) string {
	if strings.HasPrefix(coverUrl, "http") {
		return coverUrl
	}
	switch position {
	case "local":
		return fmt.Sprintf("%s/%s", svcCtx.StorageBaseUrl.Local, coverUrl)
	default:
		return coverUrl
	}
}

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
