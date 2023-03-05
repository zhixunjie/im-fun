package protocol

// 各类的操作：握手、授权认证、消息发送、消息接收

type Operation int32

const (
	OpHandshake       = Operation(iota) // handshake
	OpHandshakeReply                    //
	OpHeartbeat                         // heartbeat
	OpHeartbeatReply                    //
	OpSendMsg                           // send message.
	OpSendMsgReply                      //
	OpDisconnectReply                   // disconnect reply
	OpAuth                              // auth connect
	OpAuthReply                         //
	OpBatchMsg                          // batch message
	OpProtoReady                        // proto ready
	OpProtoFinish                       // proto finish
	OpChangeRoom                        // change room
	OpChangeRoomReply                   //
	OpSub                               // subscribe message
	OpSubReply                          //
	OpUnsub                             // unsubscribe message
	OpUnsubReply                        //
)
