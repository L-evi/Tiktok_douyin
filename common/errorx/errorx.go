package errorx

// 用于定义全局错误码
// 以及当 Gateway 请求 rpc模块 时，如果 rpc 模块返回错误，将会将错误信息转换为 ErrorResp

import (
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidParameter         = status.Error(1001, "参数无效")
	ErrFrequentVisits           = status.Error(1003, "访问频繁")
	ErrNumberInterfaceCallsOver = status.Error(1004, "接口调用次数已达上限")
	ErrTokenExperition          = status.Error(2001, "用户未登录，无权限或当前令牌已过期")
	ErrLoginTimeout             = status.Error(2002, "登录超时")
	ErrIpBlock                  = status.Error(2003, "IP被禁止访问")
	ErrLoginError               = status.Error(2004, "帐号或密码错误")
	ErrAccountError             = status.Error(2005, "帐号状态异常")
	ErrServerError              = status.Error(3001, "服务器错误")
	ErrDatabaseError            = status.Error(3002, "数据库错误")
	ErrIoOperationError         = status.Error(3003, "IO操作错误")
	ErrSystemError              = status.Error(3004, "系统错误")
	ErrRequestTimeout           = status.Error(3005, "请求超时")
)

type ErrorResp struct {
	Code int32  `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// RespErrFormat 格式化错误响应输出
func RespErrFormat(Code int32, Msg string) *ErrorResp {
	return &ErrorResp{
		Code: Code,
		Msg:  Msg,
	}
}

// FromRpcStatus 格式化 status.Error 为标准错误响应输出
func FromRpcStatus(err error) ErrorResp {
	info, ok := status.FromError(err)
	if !ok {
		return ErrorResp{
			Code: 3004,
			Msg:  "系统错误",
		}
	}

	return ErrorResp{
		Code: int32(info.Code()),
		Msg:  info.Message(),
	}
}
