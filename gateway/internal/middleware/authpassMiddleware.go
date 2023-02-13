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

type AuthPassMiddleware struct {
	identityRpc identityclient.Identity
}

func NewAuthPassMiddleware(IdentityRpcConf zrpc.RpcClientConf) *AuthMiddleware {
	_identityService := identityclient.NewIdentity(zrpc.MustNewClient(IdentityRpcConf))
	// logx.Info("identityRpc: ", _identityService)

	return &AuthMiddleware{
		identityRpc: _identityService,
	}
}

// Handle 用于处理用户认证
// token 无效 继续向下传递, 单会携带一个 is_login 为 false 的 context
// 有效则携带 User_id, username 以及一个 is_login 为 true 的 context
func (m *AuthPassMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token in query
		var token string

		if token = r.URL.Query().Get("token"); token == "" {
			token = r.PostFormValue("token")
		}

		logx.Info("token: ", token)

		// verify token in identity service
		var resp *identityclient.StatusResp
		var err error
		if resp, err = m.identityRpc.Status(r.Context(), &identityclient.StatusReq{
			Token: token,
		}); !errorx.IsRpcError(err, errorx.ErrTokenInvalid) && err != nil {
			logx.WithContext(r.Context()).Errorf("identityRpc.Status error: %v", err)
			httpx.WriteJson(w, http.StatusInternalServerError, errorx.FromRpcStatus(err))
			return
		}

		// 传递 User_id
		ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
		ctx = context.WithValue(ctx, "username", resp.Username)
		ctx = context.WithValue(ctx, "is_login", !errorx.IsRpcError(err, errorx.ErrTokenInvalid))
		next(w, r.WithContext(ctx))
	}
}
