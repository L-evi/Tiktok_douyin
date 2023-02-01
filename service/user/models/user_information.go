package models

type User_information struct {
	ID             int64  `gorm:"primary_key;auto_increment" json:"id"`
	Username       string `gorm:"size:32;not null;unique" json:"username"`
	Name           string `gorm:"size:128" json:"name"`
	Follow_count   int64  `json:"followCount"`
	Follower_count int64  `json:"followerCount"`
	Is_follow      bool   `json:"isFollow"`
}
