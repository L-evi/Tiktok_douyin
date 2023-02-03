package errx

import (
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/internal/types"
)

var (
	SUCCESS_RESP = types.Resp{Code: 0, Msg: "success"}
)

func HandleRpcErr(err error) types.Resp {
	_, ok := status.FromError(err)
	if !ok {
		logx.Errorf("error when handler err: %v", err)
		return types.Resp{
			Code: 3004,
			Msg:  err.Error(),
		}
	}

	parse := errorx.FromRpcStatus(err)

	return types.Resp{
		Code: parse.Code,
		Msg:  parse.Msg,
	}
}
