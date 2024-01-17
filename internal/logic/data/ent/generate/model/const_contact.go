package model

// 联系人状态
const (
	ContactStatusNormal = 0 // 正常
	ContactStatusDel    = 1 // 删除
)

// 联系人类型
// 0-99业务自己扩展，100之后保留
const (
	PeerNotExist = -1
	PeerNormal   = 0   // 普通用户（peer_id等于用户id）
	PeerSys      = 100 // 系统用户（peer_id等于100000）
	PeerGroup    = 101 // 群组（peer_id等于群组id）
)

// 是否给owner发过消息
const (
	PeerNotAck = 0 // 未发过
	PeerAck    = 1 // 发过
)
