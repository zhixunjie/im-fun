package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) Login(ctx *gin.Context) {
	// request
	var req request.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzUser.Login(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) RefreshToken(ctx *gin.Context) {
	// request
	var req request.RefreshTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// biz
	resp, err := s.BzUser.RefreshToken(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "error: " + err.Error()})
		return
	}

	// resp
	ctx.JSON(http.StatusOK, resp)
	return
}
