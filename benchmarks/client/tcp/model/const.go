package model

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"time"
)

const (
	OpHeartbeat      = int32(protocol.OpHeartbeat)
	OpHeartbeatReply = int32(protocol.OpHeartbeatReply)
	OpAuth           = int32(protocol.OpAuth)
	OpAuthReply      = int32(protocol.OpAuthReply)
	OpBatchMsg       = int32(protocol.OpBatchMsg)
)

const (
	RawHeaderLen = uint16(16)
	Heart        = 240 * time.Second
)
