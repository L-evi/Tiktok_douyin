package models

type Chat struct {
	ID         int64  `gorm:"primary_key;auto_increment" json:"id"`
	FromUserId int64  `gorm:"index;"`
	ToUserId   int64  `gorm:"index;"`
	Content    string `gorm:"type:text;"`
	CreateAt   int64  `gorm:"index;autoUpdateTime:milli;"`
}
