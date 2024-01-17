package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) sendToUserKeys(ctx *gin.Context) {
	// request
	var req request.SendToUserKeysReq
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
	err := s.bz.SendToUserKeys(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserKeysResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) sendToUserIds(ctx *gin.Context) {
	// request
	var req request.SendToUserIdsReq
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
	err := s.bz.SendToUserIds(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserIdsResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) sendToRoom(ctx *gin.Context) {
	// request
	var req request.SendToRoomReq
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
	err := s.bz.SendToRoom(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserRoomResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) sendToAll(ctx *gin.Context) {
	// request
	var req request.SendToAllReq
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
	err := s.bz.SendToAll(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.PushUserAllResp
	ctx.JSON(http.StatusOK, resp)
	return
}
