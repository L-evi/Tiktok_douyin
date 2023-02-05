package models

// Video 用户视频列表
type Video struct {
	ID            int64  `gorm:"primary_key;auto_increment" json:"id"`
	UserID        int64  `gorm:"index;"`
	Title         string `gorm:"size:255"`
	PlayUrl       string `gorm:"size:255"`
	CoverUrl      string `gorm:"type:text"`
	FavoriteCount int64  `gorm:"size:11"`
	CommentCount  int64  `gorm:"size:11"`
	Position      string `gorm:"size:10"` // video 存储节点 cos / local
	CreateAt      int64  `gorm:"index;autoCreateTime"`
}

// Redis SET
// KEY: prefix::video_like::[VideoId]
// VALUE: uint64
