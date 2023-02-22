package handler

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"train-tiktok/gateway/internal/types"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, types.Resp{
			Code: 404,
			Msg:  fmt.Sprintf("Route %s not exists", r.URL.Path),
		})
	})
}
