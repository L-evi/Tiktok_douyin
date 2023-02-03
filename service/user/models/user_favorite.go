package models

// UserFavorite 用户点赞视频列表
type UserFavorite struct {
	ID       int64 `gorm:"primary_key;auto_increment" json:"id"`
	UserId   int64
	VideoId  int64
	CreateAt int
}

// redis ZSET
// KEY: prefix::like_list:[UserID]
// VALUE: VideoId timeUnix
