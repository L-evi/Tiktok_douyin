package models

// VideoStat 视频点赞总数统计
type VideoStat struct {
	ID        int64
	VideoId   int
	LikeCount uint64
}

// Redis SET
// KEY: prefix::video_stat::[VideoId]
// VALUE: uint64
