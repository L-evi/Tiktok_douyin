package errorx

// 用于定义全局错误码
// 以及当 Gateway 请求 rpc模块 时，如果 rpc 模块返回错误，将会将错误信息转换为 ErrorResp

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidParameter         = status.Error(1001, "参数无效")
	ErrFrequentVisits           = status.Error(1003, "访问频繁")
	ErrNumberInterfaceCallsOver = status.Error(1004, "接口调用次数已达上限")
	ErrTokenInvalid             = status.Error(2001, "用户未登录，无权限或当前令牌已过期")
	ErrLoginTimeout             = status.Error(2002, "登录超时")
	ErrLoginError               = status.Error(2004, "帐号或密码错误")
	ErrDatabaseError            = status.Error(3001, "数据库错误")
	ErrIoOperationError         = status.Error(3002, "IO操作错误")
	ErrSystemError              = status.Error(3003, "服务器开小差了, 过会儿再试吧")
	ErrRequestTimeout           = status.Error(3004, "请求超时")
	ErrFileTypeNotSupport       = status.Error(4001, "文件类型不支持")
	ErrUserNotFound             = status.Error(5001, "用户不存在")
	ErrNotFriend                = status.Error(5002, "不是好友关系")
)

type ErrorResp struct {
	Code int32  `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// RespErrFormat 格式化错误响应输出
func RespErrFormat(Code int32, Msg string) ErrorResp {
	return ErrorResp{
		Code: Code,
		Msg:  Msg,
	}
}

func NewError(code int32, msg string) ErrorResp {
	return ErrorResp{
		Code: code,
		Msg:  msg,
	}
}

func ParseRpcError(err error) error {
	info, ok := status.FromError(err)
	if !ok {
		return ErrSystemError
	}
	return info.Err()
}

func IsRpcError(rpcErr error, err error) bool {
	if errors.Is(ParseRpcError(rpcErr), err) {
		return true
	}
	return false
}

// FromRpcStatus 格式化 status.Error 为标准错误响应输出
func FromRpcStatus(err error) ErrorResp {
	info, ok := status.FromError(err)
	if !ok {
		return FromRpcStatus(ErrSystemError)
	}

	// 非自定义错误, 禁止向用户返回详细错误信息
	if info.Code() < 100 {
		logx.Errorf("from rpc status error: %v", err)
		return FromRpcStatus(ErrSystemError)
	}

	return RespErrFormat(int32(info.Code()), info.Message())
}
