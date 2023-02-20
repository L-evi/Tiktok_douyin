package errx

import "google.golang.org/grpc/status"

var (
	ErrCantSendToSelf = status.Error(1004001, "can't send message to self")
)
