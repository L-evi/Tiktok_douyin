package errx

import "google.golang.org/grpc/status"

var (
	ErrRepeatFollow = status.Error(1003001, "重复关注")
)
