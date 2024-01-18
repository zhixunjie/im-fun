package model

// 联系人状态
const (
	ContactStatusNormal = 0 // 正常
	ContactStatusDel    = 1 // 删除
)

// 联系人类型
// 0-99业务自己扩展，100之后保留
const (
	PeerTypeNormal = 0   // 普通用户（peer_id等于用户id）
	PeerTypeSys    = 100 // 系统用户（peer_id等于100000）
	PeerTypeGroup  = 101 // 群组（peer_id等于群组id）
)

// 是否给owner发过消息
const (
	PeerNotAck = 0 // 未发过
	PeerAck    = 1 // 发过
)

type BuildContactParams struct {
	MsgId    uint64
	OwnerId  uint64
	PeerId   uint64
	PeerType int32
	PeerAck  uint32
}
