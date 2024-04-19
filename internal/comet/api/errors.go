package api

import (
	"errors"
)

var (
	// ErrHandshake server
	ErrHandshake = errors.New("handshake failed")
	ErrOperation = errors.New("request operation not valid")

	// ErrTimerFull timer
	ErrTimerFull   = errors.New("timer full")
	ErrTimerEmpty  = errors.New("timer empty")
	ErrTimerNoItem = errors.New("timer item not exist")

	// ErrParamsNotAllow channel
	ErrParamsNotAllow    = errors.New("params error")
	ErrChannelSignalFull = errors.New("signal channel full and msg will be dropped")

	ErrBroadCastRoomArg = errors.New("rpc broadcast room arg error")

	// ErrRoomDrop room
	ErrRoomDrop = errors.New("room has drop")

	// ErrRPCLogic rpc
	ErrRPCLogic = errors.New("logic rpc is not available")

	ErrTCPWriteError = errors.New("write err")
)
