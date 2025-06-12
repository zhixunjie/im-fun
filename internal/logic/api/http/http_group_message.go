package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) GroupMessageSend(ctx *gin.Context) {
	// request
	var req request.GroupMessageSendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzGroupMessage.Send(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) GroupMessageFetch(ctx *gin.Context) {
	// request
	var req request.GroupMessageFetchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzGroupMessage.Fetch(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}
