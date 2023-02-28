package apihttp

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"net/http"
)

func (s *Server) pingHandler(ctx *gin.Context) {
	// request
	var req request.PingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// service
	s.svc.Ping()

	// resp
	var resp response.PingResp
	resp.Pong = "pong"
	ctx.JSON(http.StatusOK, resp)
	return
}
