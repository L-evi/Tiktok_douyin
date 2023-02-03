package models

// VideoLike 视频点赞总数持久化
type VideoLike struct {
	ID        int32 `gorm:"primary_key;auto_increment" json:"id"`
	VideoId   int64
	LikeCount uint64
}

// Redis SET
// KEY: prefix::video_like::[VideoId]
// VALUE: uint64
