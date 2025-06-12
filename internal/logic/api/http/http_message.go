package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) MessageSend(ctx *gin.Context) {
	// request
	var req request.MessageSendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzMessage.Send(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) MessageFetch(ctx *gin.Context) {
	// request
	var req request.MessageFetchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzMessage.Fetch(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) MessageClearHistory(ctx *gin.Context) {
	// request
	var req request.MessageClearHistoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzMessage.ClearHistory(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}
