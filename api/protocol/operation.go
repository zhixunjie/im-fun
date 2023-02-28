package protocol

// 各类的操作：握手、

type Operation int32

const (
	// OpHandshake handshake
	OpHandshake = Operation(iota)
	OpHandshakeReply
	// OpHeartbeat heartbeat
	OpHeartbeat
	OpHeartbeatReply
	// OpSendMsg send message.
	OpSendMsg
	OpSendMsgReply
	// OpDisconnectReply disconnect reply
	OpDisconnectReply
	// OpAuth auth connect
	OpAuth
	OpAuthReply
	// OpRaw raw message
	OpRaw
	// OpProtoReady proto ready
	OpProtoReady
	// OpProtoFinish proto finish
	OpProtoFinish
	// OpChangeRoom change room
	OpChangeRoom
	OpChangeRoomReply
	// OpSub subscribe operation
	OpSub
	OpSubReply
	// OpUnsub unsubscribe operation
	OpUnsub
	OpUnsubReply
)
