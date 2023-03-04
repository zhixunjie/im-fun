package apihttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"net/http"
)

func (s *Server) pushUserKeys(ctx *gin.Context) {
	// request
	var req request.PushUserKeysReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// check params
	if len(req.UserKeys) == 0 {
		response.JsonError(ctx, errors.New("req.UserKeys not allow"))
		return
	}
	if len(req.Message) == 0 {
		response.JsonError(ctx, errors.New("req.Message not allow"))
		return
	}

	// service
	err := s.svc.PushUserKeys(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserKeysResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) pushUserIds(ctx *gin.Context) {
	// request
	var req request.PushUserIdsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// check params
	if len(req.UserIds) == 0 {
		response.JsonError(ctx, errors.New("req.UserIds not allow"))
		return
	}
	if len(req.Message) == 0 {
		response.JsonError(ctx, errors.New("req.Message not allow"))
		return
	}

	// service
	err := s.svc.PushUserIds(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserIdsResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) pushUserRoom(ctx *gin.Context) {
	// request
	var req request.PushUserRoomReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// check params
	if len(req.RoomId) == 0 {
		response.JsonError(ctx, errors.New("req.RoomId not allow"))
		return
	}
	if len(req.RoomType) == 0 {
		response.JsonError(ctx, errors.New("req.RoomType not allow"))
		return
	}
	if len(req.Message) == 0 {
		response.JsonError(ctx, errors.New("req.Message not allow"))
		return
	}

	// service
	err := s.svc.PushUserRoom(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserRoomResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) pushUserAll(ctx *gin.Context) {
	// request
	var req request.PushUserAllReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// check params
	if req.Speed == 0 {
		response.JsonError(ctx, errors.New("req.Speed not allow"))
		return
	}
	if len(req.Message) == 0 {
		response.JsonError(ctx, errors.New("req.Message not allow"))
		return
	}

	// service
	err := s.svc.PushUserAll(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserAllResp
	ctx.JSON(http.StatusOK, resp)
	return
}
