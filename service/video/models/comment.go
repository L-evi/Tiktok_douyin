package models

type Comment struct {
	ID         int64  `gorm:"primary_key;auto_increment" json:"id"`
	VideoID    int64  `gorm:"index;"`
	UserID     int64  `gorm:"index;"`
	Content    string `gorm:"type:text"`
	CreateDate string `gorm:"index;autoCreateTime"`
}
