package protocol

// 各类的操作：握手、授权认证、消息发送、消息接收

// Operation 定义所有 WebSocket 通信操作码
type Operation int32

const (
	// ===== 握手 Handshake =====
	OpHandshake      Operation = 0 // 客户端请求握手
	OpHandshakeReply Operation = 1 // 服务端响应握手

	// ===== 心跳 Heartbeat =====
	OpHeartbeat     Operation = 2 // 客户端发送心跳
	OpHeartbeatResp Operation = 3 // 服务端回应心跳

	// ===== 消息 Message =====
	OpSendMsg         Operation = 4 // 客户端发送消息
	OpSendMsgReply    Operation = 5 // 服务端回应消息
	OpDisconnectReply Operation = 6 // 服务端断开连接响应

	// ===== 鉴权 Auth =====
	OpAuth     Operation = 7 // 客户端发送认证
	OpAuthResp Operation = 8 // 服务端认证响应

	// ===== 批量消息 Batch / 原始消息 Raw =====
	OpBatchMsg Operation = 9 // 批量或原始消息下发

	// ===== 协议状态 Protocol Status =====
	OpProtoReady  Operation = 10 // 协议准备完毕
	OpProtoFinish Operation = 11 // 协议处理完成

	// ===== 房间 Room =====
	OpChangeRoom     Operation = 12 // 切换房间
	OpChangeRoomResp Operation = 13 // 切换房间响应

	// ===== 订阅 Subscribe =====
	OpSub      Operation = 14 // 订阅消息
	OpSubReply Operation = 15 // 订阅响应

	// ===== 取消订阅 Unsubscribe =====
	OpUnsub      Operation = 16 // 取消订阅消息
	OpUnsubReply Operation = 17 // 取消订阅响应
)
