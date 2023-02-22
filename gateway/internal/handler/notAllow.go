package handler

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"train-tiktok/gateway/internal/types"
)

func NotAllowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJsonCtx(r.Context(), w, http.StatusMethodNotAllowed, types.Resp{
			Code: 405,
			Msg:  fmt.Sprintf("Method %s not allowed for route %s", r.Method, r.URL.Path),
		})
	})
}
