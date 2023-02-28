package protocol

type Operation int32

const (
	// OpHandshake handshake
	OpHandshake = Operation(iota)
	// OpHandshakeReply handshake reply
	OpHandshakeReply

	// OpHeartbeat heartbeat
	OpHeartbeat
	// OpHeartbeatReply heartbeat reply
	OpHeartbeatReply

	// OpSendMsg send message.
	OpSendMsg
	// OpSendMsgReply  send message reply
	OpSendMsgReply

	// OpDisconnectReply disconnect reply
	OpDisconnectReply

	OpAuth // OpAuth auth connect
	// OpAuthReply auth connect reply
	OpAuthReply

	// OpRaw raw message
	OpRaw

	// OpProtoReady proto ready
	OpProtoReady
	// OpProtoFinish proto finish
	OpProtoFinish

	// OpChangeRoom change room
	OpChangeRoom
	// OpChangeRoomReply change room reply
	OpChangeRoomReply

	// OpSub subscribe operation
	OpSub
	// OpSubReply subscribe operation
	OpSubReply

	// OpUnsub unsubscribe operation
	OpUnsub
	// OpUnsubReply unsubscribe operation reply
	OpUnsubReply
)
