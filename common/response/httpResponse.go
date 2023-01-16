package response

type HttpResponse struct {
	responseStatus
	Data interface{} `json:"data"`
}

func (res *HttpResponse) WithMsg(msg string) HttpResponse {
	return HttpResponse{
		responseStatus: responseStatus{
			Code: res.responseStatus.Code,
			Msg:  msg,
		},
		Data: res.Data,
	}
}

func (res *HttpResponse) WithData(data interface{}) HttpResponse {
	return HttpResponse{
		responseStatus: responseStatus{
			Code: res.responseStatus.Code,
			Msg:  res.responseStatus.Msg,
		},
		Data: data,
	}
}

func Response(status responseStatus) *HttpResponse {
	return &HttpResponse{
		responseStatus: responseStatus{
			Code: status.Code,
			Msg:  status.Msg,
		},
		Data: nil,
	}
}
