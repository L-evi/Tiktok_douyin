package test

import (
	"Tiktok_douyin/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestHttpResponse(t *testing.T) {
	engine := gin.Default()
	res := struct {
		Name string `json:"name"`
	}{
		Name: "one",
	}

	engine.GET("/ok", func(context *gin.Context) {
		context.JSON(http.StatusOK, response.Response(response.SUCCESS).WithData(res))
	})
	engine.Run(":8080")
}
