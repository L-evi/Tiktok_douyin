package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/router"
	"google.golang.org/grpc/status"
	"net/http"
	"os/exec"
	"strings"
	"train-tiktok/common/errorx"
	"train-tiktok/gateway/internal/config"
	"train-tiktok/gateway/internal/handler"
	"train-tiktok/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

// videoStaticHandler 视频对外路由
func videoStaticHandler(engine *rest.Server, conf config.Config) {
	level := []string{":1", ":2", ":3", ":4"}
	pattern := fmt.Sprintf("/%s/", conf.PublicPath)
	fileDir := fmt.Sprintf("./%s/", conf.PublicPath)
	for i := 1; i < len(level); i++ {
		path := "/" + strings.Join(level[:i], "/")
		//最后生成 /asset
		engine.AddRoute(
			rest.Route{
				Method: http.MethodGet,
				Path:   path,
				Handler: func(pattern string, fileDir string) http.HandlerFunc {
					// log.Println(fileDir, pattern)
					return func(w http.ResponseWriter, req *http.Request) {
						http.StripPrefix(pattern, http.FileServer(http.Dir(fileDir))).ServeHTTP(w, req)
					}
				}(pattern, fileDir),
			})
		//log.Println("videoStaticHandler", path, pattern, fileDir)
	}
}

func main() {
	flag.Parse()

	// check if ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		panic("ffmpeg is not installed")
	}

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	c = ctx.Config

	// set router
	r := router.NewRouter()
	r.SetNotAllowedHandler(handler.NotAllowHandler())
	r.SetNotFoundHandler(handler.NotFoundHandler())

	server := rest.MustNewServer(c.RestConf, rest.WithRouter(r))
	defer server.Stop()

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		_, ok := status.FromError(err)
		if !ok {
			logx.WithContext(ctx).Errorf("error when handler rpc err: %v", err)
			return http.StatusOK, errorx.ErrorResp{Code: 3004, Msg: err.Error()}
		}

		// logx.WithContext(ctx).Errorf("error when handler err: %v", err)
		return http.StatusOK, errorx.FromRpcStatus(err)
	})

	handler.RegisterHandlers(server, ctx)
	// 注册视频对外路由
	videoStaticHandler(server, c)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
