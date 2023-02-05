package models

// Video 用户视频列表
type Video struct {
	ID            int64  `gorm:"primary_key;auto_increment" json:"id"`
	UserID        int64  `gorm:"index;"`
	Title         string `gorm:"size:255"`
	PlayUrl       string `gorm:"size:255"`
	CoverUrl      string `gorm:"type:text"`
	FavoriteCount int64
	CommentCount  int64
}

// Redis SET
// KEY: prefix::video_like::[VideoId]
// VALUE: uint64
