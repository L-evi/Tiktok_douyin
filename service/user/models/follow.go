package models

// Follow 用户关注表
type Follow struct {
	ID       int64 `gorm:"primary_key;auto_increment" json:"id"`
	UserId   int64 `gorm:"index"`
	TargetId int64 `gorm:"index"`
	CreateAt int64 `gorm:"autoCreateTime"`
}
