package errors

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

	// ErrPushMsgArg channel
	ErrPushMsgArg           = errors.New("rpc pushmsg arg error")
	ErrPushMsgsArg          = errors.New("rpc pushmsgs arg error")
	ErrMPushMsgArg          = errors.New("rpc mpushmsg arg error")
	ErrMPushMsgsArg         = errors.New("rpc mpushmsgs arg error")
	ErrSignalFullMsgDropped = errors.New("signal channel full, msg dropped")

	// ErrBroadCastArg bucket
	ErrBroadCastArg     = errors.New("rpc broadcast arg error")
	ErrBroadCastRoomArg = errors.New("rpc broadcast  room arg error")

	// ErrRoomDrop room
	ErrRoomDrop = errors.New("room has drop")

	// ErrRPCLogic rpc
	ErrRPCLogic = errors.New("logic rpc is not available")
)
