package errx

import (
	"google.golang.org/grpc/status"
)

var (
	ErrNoLatestVideo    = status.Error(1002101, "没有更新的视频了")
	ErrVideoNotFound    = status.Error(1002103, "没有找到视频")
	ErrAlreadyFavorite  = status.Error(1002104, "已经点赞过了")
	ErrCommentTextEmpty = status.Error(1002105, "评论内容不能为空")
)
