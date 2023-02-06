package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/identityclient"
)

type AuthMiddleware struct {
	identityRpc identityclient.Identity
}

func NewAuthMiddleware(IdentityRpcConf zrpc.RpcClientConf) *AuthMiddleware {
	_identityService := identityclient.NewIdentity(zrpc.MustNewClient(IdentityRpcConf))
	logx.Info("identityRpc: ", _identityService)

	return &AuthMiddleware{
		identityRpc: _identityService,
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
				httpx.WriteJson(w, http.StatusForbidden, errorx.FromRpcStatus(errorx.ErrTokenInvalid))
				return
			}
		}

		logx.Info("token: ", token)

		// verify token in identity service
		var resp *identityclient.StatusResp
		var err error
		if resp, err = m.identityRpc.Status(r.Context(), &identityclient.StatusReq{
			Token: token,
		}); errorx.IsRpcError(err, errorx.ErrTokenInvalid) {
			httpx.WriteJson(w, http.StatusUnauthorized, errorx.FromRpcStatus(errorx.ErrTokenInvalid))
			return
		} else if err != nil {
			logx.WithContext(r.Context()).Errorf("identityRpc.Status error: %v", err)
			httpx.WriteJson(w, http.StatusInternalServerError, errorx.FromRpcStatus(err))
		}

		// 传递 User_id
		ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
		ctx = context.WithValue(ctx, "username", resp.Username)
		next(w, r.WithContext(ctx))
	}
}
