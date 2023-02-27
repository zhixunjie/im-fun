package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	service2 "github.com/zhixunjie/im-fun/internal/logic/service"
	"net/http"
)

func sendHandler(ctx *gin.Context) {
	// request
	var req request.SendMsgReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	if req.SendId == 0 || req.ReceiveId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + "id not allow"})
		return
	}

	// service
	resp, err := service2.SendMessage(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func fetchHandler(ctx *gin.Context) {
	// request
	var req request.PingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// service
	service2.Ping()

	// resp
	var resp response.PingResp
	resp.Pong = "pong"
	ctx.JSON(http.StatusOK, resp)
	return
}
