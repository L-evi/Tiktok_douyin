package response

type httpResponse struct {
	responseStatus
	Data interface{} `json:"data"`
}

func (res *httpResponse) WithMsg(msg string) httpResponse {
	return httpResponse{
		responseStatus: responseStatus{
			Code: res.responseStatus.Code,
			Msg:  msg,
		},
		Data: res.Data,
	}
}

func (res *httpResponse) WithData(data interface{}) httpResponse {
	return httpResponse{
		responseStatus: responseStatus{
			Code: res.responseStatus.Code,
			Msg:  res.responseStatus.Msg,
		},
		Data: data,
	}
}

func HttpResponse(status responseStatus) *httpResponse {
	return &httpResponse{
		responseStatus: responseStatus{
			Code: status.Code,
			Msg:  status.Msg,
		},
		Data: nil,
	}
}
