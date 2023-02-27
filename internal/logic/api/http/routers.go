package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"net/http"
)

// SetupRouter 设置-路由
func SetupRouter(router *gin.Engine) {
	// 设置-单个路由
	router.GET("/ping", pingHandler)

	// 设置-路由组
	g1 := router.Group("/message")
	{
		g1.POST("/send", sendHandler)
		g1.GET("/fetch", pingHandler)
	}
}

func pingHandler(ctx *gin.Context) {
	// request
	var req request.PingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// service
	service.Ping()

	// resp
	var resp response.PingResp
	resp.Pong = "pong"
	ctx.JSON(http.StatusOK, resp)
	return
}
