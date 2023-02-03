package models

type User struct {
	ID       int64  `gorm:"primary_key;auto_increment" json:"id"`
	Username string `gorm:"size:32;not null;unique"`
	Nickname string `gorm:"size:128;not null"`
	Password string `gorm:"size:128;not null"`
}
