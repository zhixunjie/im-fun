package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) sendMessage(ctx *gin.Context) {
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
	resp, err := s.BzMessage.Send(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) fetchMessage(ctx *gin.Context) {
	// request
	var req request.FetchMsgReq
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

func (s *Server) fetchContact(ctx *gin.Context) {
	// request
	var req request.FetchContactReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzContact.FetchContact(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}
