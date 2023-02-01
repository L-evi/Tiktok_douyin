package main

import (
	"Tiktok_douyin/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()
	res := struct {
		Name string `json:"name"`
	}{
		Name: "one",
	}

	engine.GET("/ok", func(context *gin.Context) {
		context.JSON(http.StatusOK, response.HttpResponse(response.SUCCESS).WithData(res))
	})
	engine.Run(":8080")
}
