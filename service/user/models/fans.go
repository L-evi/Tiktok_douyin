package models

// Fans 用户粉丝表
type Fans struct {
	ID       int64 `gorm:"primary_key;auto_increment" json:"id"`
	UserId   int64 `gorm:"index"`
	TargetId int64
}
