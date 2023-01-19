package errutil

import (
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUsername = status.Error(1001101, "用户名应为5-32位的数字或字母")
	ErrInvalidPassword = status.Error(1001102, "密码应为5-32位的字符")

	ErrUsernameExists = status.Error(1001103, "用户名已存在")

	ErrWrongIdentity = status.Error(1001104, "用户名或密码错误")
	ErrTokenInvalid  = status.Error(1001105, "用户登录失效, 请重新登录")
)
