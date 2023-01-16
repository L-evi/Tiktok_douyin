package response

type responseStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	SUCCESS = responseStatus{
		Code: 200,
		Msg:  "成功",
	}

	INVALID_PARAMETER = responseStatus{
		Code: 1001,
		Msg:  "参数无效",
	}
	FREQUENT_VISITS = responseStatus{
		Code: 1003,
		Msg:  "访问频繁",
	}
	NUMBER_INTERFACE_CALLS_OVER = responseStatus{
		Code: 1004,
		Msg:  "接口调用次数已达上限",
	}
	TOKEN_EXPERITION = responseStatus{
		Code: 2001,
		Msg:  "用户未登录，无权限或当前令牌已过期",
	}
	LOGIN_TIMEOUT = responseStatus{
		Code: 2002,
		Msg:  "登录超时",
	}
	IP_BLOCK = responseStatus{
		Code: 2003,
		Msg:  "IP被禁止访问",
	}
	LOGIN_ERROR = responseStatus{
		Code: 2004,
		Msg:  "帐号或密码错误",
	}
	ACCOUNT_ERROR = responseStatus{
		Code: 2005,
		Msg:  "帐号状态异常",
	}
	SERVER_ERROR = responseStatus{
		Code: 3001,
		Msg:  "服务器错误",
	}

	DATABASE_ERROR = responseStatus{
		Code: 3002,
		Msg:  "数据库错误",
	}

	IO_OPERATION_ERROR = responseStatus{
		Code: 3003,
		Msg:  "IO操作错误",
	}

	SYSTEM_ERROR = responseStatus{
		Code: 3004,
		Msg:  "系统错误",
	}

	REQUEST_TIMEOUT = responseStatus{
		Code: 3005,
		Msg:  "请求超时",
	}
)
