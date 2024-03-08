package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"net/http"
)

func (s *Server) sendToUsers(ctx *gin.Context) {
	// request
	var req request.SendToUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}

	// check params
	if len(req.TcpSessionIds) == 0 {
		err := errors.New("SessionIds not allow")
		response.JsonError(ctx, err)
		return
	}
	if len(req.Message) == 0 {
		err := errors.New("message not allow")
		response.JsonError(ctx, err)
		return
	}

	// invoke service
	err := s.bz.SendToUsers(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.SendToUsersResp
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Server) sendToUsersByIds(ctx *gin.Context) {
	// request
	var req request.SendToUsersByIdsReq
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

	// invoke service
	err := s.bz.SendToUsersByIds(ctx, &req)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// resp
	var resp response.SendToUsersByIdsResp
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
	var resp response.SendToRoomResp
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
	var resp response.SendToAllResp
	ctx.JSON(http.StatusOK, resp)
	return
}
