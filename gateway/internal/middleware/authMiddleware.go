package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
	"net/http"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/identityclient"
)

type AuthMiddleware struct {
	IdentityRpcConf zrpc.RpcClientConf
}

func NewAuthMiddleware(IdentityRpcConf zrpc.RpcClientConf) *AuthMiddleware {
	return &AuthMiddleware{
		IdentityRpcConf: IdentityRpcConf,
	}
}

// Handle 用于处理用户认证
// token 无效则直接返回 403 / 登录超时
// 有效则携带 UserId 以及 Username 向后传递
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token in query
		var token string

		if token = r.URL.Query().Get("token"); token == "" {
			if token = r.PostFormValue("token"); token == "" {
				httpx.WriteJson(w, http.StatusUnauthorized, errorx.ErrLoginTimeout)
				return
			}
		}

		log.Println("token: ", token)

		// verify token in identity service
		identityRpc := identityclient.NewIdentity(zrpc.MustNewClient(m.IdentityRpcConf))

		logx.Info("identityRpc: ", identityRpc)

		var resp *identityclient.StatusResp
		var err error
		if resp, err = identityRpc.Status(r.Context(), &identityclient.StatusReq{
			Token: token,
		}); err != nil {
			logx.Errorf("Auth middleware: identity rpc status err: %v", err)
			httpx.WriteJson(w, http.StatusInternalServerError, errorx.ErrSystemError)
			return
		}

		// process request from unlogged user
		if resp.IsLogin == false {
			httpx.WriteJson(w, http.StatusUnauthorized, errorx.ErrLoginTimeout)
			return
		}

		// 传递 User_id
		ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
		ctx = context.WithValue(ctx, "username", resp.Username)
		next(w, r.WithContext(ctx))
	}
}
