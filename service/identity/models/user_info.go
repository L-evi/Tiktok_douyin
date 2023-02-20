package models

// UserInformation 用户详细信息表
type UserInformation struct {
	ID              int64 `gorm:"primary;auto_increment" json:"id"`
	UserId          int64
	Nickname        string
	Avatar          string
	BackgroundImage string
	Signature       string
}
