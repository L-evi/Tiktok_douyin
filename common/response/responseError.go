package response

import "fmt"

type Error interface {
	error
	Error() string
	GetCode() int32
}

type RespError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (e *RespError) Error() string {
	return fmt.Sprintf("err_code: %d, err_msg: %s", e.Code, e.Msg)
}

func (e *RespError) GetCode() int32 {
	return e.Code
}

// NewErr creates a new error.
func NewErr(code int32, msg string) *RespError {
	return &RespError{
		Code: code,
		Msg:  msg,
	}
}

var (
	ErrInvalidParameter         = NewErr(1001, "参数无效")
	ErrFrequentVisits           = NewErr(1003, "访问频繁")
	ErrNumberInterfaceCallsOver = NewErr(1004, "接口调用次数已达上限")
	ErrTokenExperition          = NewErr(2001, "用户未登录，无权限或当前令牌已过期")
	ErrLoginTimeout             = NewErr(2002, "登录超时")
	ErrIpBlock                  = NewErr(2003, "IP被禁止访问")
	ErrLoginError               = NewErr(2004, "帐号或密码错误")
	ErrAccountError             = NewErr(2005, "帐号状态异常")
	ErrServerError              = NewErr(3001, "服务器错误")
	ErrDatabaseError            = NewErr(3002, "数据库错误")
	ErrIoOperationError         = NewErr(3003, "IO操作错误")
	ErrSystemError              = NewErr(3004, "系统错误")
	ErrRequestTimeout           = NewErr(3005, "请求超时")
)
