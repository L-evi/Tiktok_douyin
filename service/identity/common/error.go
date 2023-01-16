package common

import "train-tiktok/common/response"

var (
	ErrInvalidUsername = response.NewErr(1001101, "用户名应为5-32位的数字或字母")
	ErrInvalidPassword = response.NewErr(1001102, "密码应为5-32位的字符")

	ErrUsernameExists = response.NewErr(1001103, "用户名已存在")
)
