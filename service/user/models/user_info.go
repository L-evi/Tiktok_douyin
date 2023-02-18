package models

// UserInfomation 用户详细信息表
type UserInfo struct {
	ID              int64 `gorm:"primary;auto_increment" json:"id"`
	UserId          int64
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorite   int64
	WorkCount       int64
	FavoriteCount   int64
}
