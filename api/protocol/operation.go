package protocol

// 各类的操作：握手、授权认证、消息发送、消息接收

type Operation int32

const (
	// OpHandshake handshake
	OpHandshake = Operation(iota)
	OpHandshakeReply

	// OpHeartbeat heartbeat
	OpHeartbeat
	OpHeartbeatReply

	// OpSendMsg send message
	OpSendMsg
	OpSendMsgReply

	// OpDisconnectReply disconnect reply
	OpDisconnectReply

	// OpAuth auth connect
	OpAuth
	OpAuthReply

	// OpBatchMsg batch messages / raw messages
	OpBatchMsg

	// OpProtoReady proto ready
	OpProtoReady
	// OpProtoFinish proto finish
	OpProtoFinish

	// OpChangeRoom change room
	OpChangeRoom
	OpChangeRoomReply

	// OpSub subscribe message
	OpSub
	OpSubReply

	// OpUnsub unsubscribe message
	OpUnsub
	OpUnsubReply
)
