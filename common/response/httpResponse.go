package response

type HttpResponse struct {
	RespStatus
	Data interface{} `json:"data"`
}

func (res *HttpResponse) WithMsg(msg string) HttpResponse {
	return HttpResponse{
		RespStatus: RespStatus{
			Code: res.RespStatus.Code,
			Msg:  msg,
		},
		Data: res.Data,
	}
}

func (res *HttpResponse) WithData(data interface{}) HttpResponse {
	return HttpResponse{
		RespStatus: RespStatus{
			Code: res.RespStatus.Code,
			Msg:  res.RespStatus.Msg,
		},
		Data: data,
	}
}

func Response(status RespStatus) *HttpResponse {
	return &HttpResponse{
		RespStatus: RespStatus{
			Code: status.Code,
			Msg:  status.Msg,
		},
		Data: nil,
	}
}
