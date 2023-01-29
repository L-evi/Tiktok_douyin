package models

// UserFavorite 用户点赞视频列表
type UserFavorite struct {
	ID       int64
	UserId   int
	VideoId  int
	CreateAt int
}

// redis ZSET
// KEY: prefix::like_list:[UserID]
// VALUE: VideoId timeUnix
