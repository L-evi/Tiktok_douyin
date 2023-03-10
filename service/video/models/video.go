package models

// Video 用户视频列表
type Video struct {
	ID       int64  `gorm:"primary_key;auto_increment" json:"id"`
	UserID   int64  `gorm:"index;"`
	Title    string `gorm:"size:255"`
	PlayUrl  string `gorm:"size:255"`
	CoverUrl string `gorm:"type:text"`
	Position string `gorm:"size:10"` // video 存储节点 cos / local
	Hash     string `gorm:"index;size:255"`
	CreateAt int64  `gorm:"index;autoCreateTime"`
}

// Redis SET
// KEY: prefix::video_like::[VideoId]
// VALUE: uint64
