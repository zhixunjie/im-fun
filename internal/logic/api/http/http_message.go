package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"net/http"
)

func (s *Server) send(ctx *gin.Context) {
	// request
	var req request.SendMsgReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	if req.SendId == 0 || req.PeerId == 0 {
		response.JsonError(ctx, errors.New("id not allow"))
		return
	}

	// service
	resp, err := s.svc.SendMessage(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) fetch(ctx *gin.Context) {
	// request
	var req request.PingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
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
