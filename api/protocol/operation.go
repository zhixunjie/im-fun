package protocol

type Operation int32

const (
	// OpHandshake handshake
	OpHandshake = Operation(0)
	// OpHandshakeReply handshake reply
	OpHandshakeReply = Operation(1)

	// OpHeartbeat heartbeat
	OpHeartbeat = Operation(2)
	// OpHeartbeatReply heartbeat reply
	OpHeartbeatReply = Operation(3)

	// OpSendMsg send message.
	OpSendMsg = Operation(4)
	// OpSendMsgReply  send message reply
	OpSendMsgReply = Operation(5)

	// OpDisconnectReply disconnect reply
	OpDisconnectReply = Operation(6)

	// OpAuth auth connnect
	OpAuth = Operation(7)
	// OpAuthReply auth connect reply
	OpAuthReply = Operation(8)

	// OpRaw raw message
	OpRaw = Operation(9)

	// OpProtoReady proto ready
	OpProtoReady = Operation(10)
	// OpProtoFinish proto finish
	OpProtoFinish = Operation(11)

	// OpChangeRoom change room
	OpChangeRoom = Operation(12)
	// OpChangeRoomReply change room reply
	OpChangeRoomReply = Operation(13)

	// OpSub subscribe operation
	OpSub = Operation(14)
	// OpSubReply subscribe operation
	OpSubReply = Operation(15)

	// OpUnsub unsubscribe operation
	OpUnsub = Operation(16)
	// OpUnsubReply unsubscribe operation reply
	OpUnsubReply = Operation(17)
)
