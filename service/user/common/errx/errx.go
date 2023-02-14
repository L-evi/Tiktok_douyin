package errx

import "google.golang.org/grpc/status"

var (
	ErrRepeatFollow  = status.Error(1003001, "重复关注")
	ErrSelfFollow    = status.Error(1003002, "不能关注自己")
	ErrLoginRequired = status.Error(1003003, "请先登录")
)
