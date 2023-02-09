package errx

import (
	"google.golang.org/grpc/status"
)

var (
	ErrNoLatestVideo = status.Error(1002101, "没有更新的视频了")
)
