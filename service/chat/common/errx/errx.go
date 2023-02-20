package errx

import "google.golang.org/grpc/status"

var (
	ErrCantSendToSelf = status.Error(1004001, "不能给自己发信息")
	ErrContentEmpty   = status.Error(1004002, "内容不能为空")
)
