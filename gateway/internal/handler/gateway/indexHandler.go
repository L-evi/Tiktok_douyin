package gateway

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"train-tiktok/gateway/internal/logic/gateway"
	"train-tiktok/gateway/internal/svc"
)

func IndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := gateway.NewIndexLogic(r.Context(), svcCtx)
		resp, err := l.Index()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
