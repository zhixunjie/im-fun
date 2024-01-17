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

	// biz
	resp, err := s.bzMessage.SendMessage(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) fetch(ctx *gin.Context) {
	// request
	var req request.FetchMsgReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.bzMessage.FetchMessage(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}
