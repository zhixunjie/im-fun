package protocol

// 各类的操作：握手、授权认证、消息发送、消息接收

// Operation 定义所有 WebSocket 通信操作码
type Operation int32

const (
	// ===== 握手 Handshake =====
	OpHandshake      Operation = 0 // client: 请求握手
	OpHandshakeReply Operation = 1 // server: 响应握手

	// ===== 心跳 Heartbeat =====
	OpHeartbeatReq  Operation = 2 // client: 发送心跳
	OpHeartbeatResp Operation = 3 // server: 回应心跳

	// ===== 消息 Message =====
	OpSendMsg         Operation = 4 // client: 发送消息
	OpSendMsgReply    Operation = 5 // server: 回应消息
	OpDisconnectReply Operation = 6 // server: 断开连接响应

	// ===== 鉴权 Auth =====
	OpAuthReq  Operation = 7 // client: 发送认证
	OpAuthResp Operation = 8 // server: 认证响应

	// ===== 批量消息 Batch / 原始消息 Raw =====
	OpBatchMsg Operation = 9 // server: 批量或原始消息下发

	// ===== 协议状态 Protocol Status =====
	OpProtoReady  Operation = 10 // 协议准备完毕
	OpProtoFinish Operation = 11 // 协议处理完成

	// ===== 房间 Room =====
	OpChangeRoomReq  Operation = 12 // client: 切换房间
	OpChangeRoomResp Operation = 13 // server: 切换房间响应

	// ===== 订阅 Subscribe =====
	OpSubReq    Operation = 14 // client: 订阅消息
	OpSubResp   Operation = 15 // server: 订阅响应
	OpUnsub     Operation = 16 // client: 取消订阅消息
	OpUnsubResp Operation = 17 // server: 取消订阅响应
)
