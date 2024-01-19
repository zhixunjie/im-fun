package model

// ContactStatus 联系人状态
type ContactStatus uint32

const (
	ContactStatusNormal = 0 // 正常
	ContactStatusDel    = 1 // 删除
)

// ==============================================

// PeerType 联系人类型
// 0-99业务自己扩展，100之后保留
type PeerType int32

const (
	PeerTypeNormal PeerType = 0   // 普通用户（peer_id等于用户id）
	PeerTypeSystem PeerType = 100 // 系统用户（peer_id等于100000）
	PeerTypeGroup  PeerType = 101 // 群组（peer_id等于群组id）

	// SystemUid 系统用户的用户ID
	SystemUid = 100000
)

// ==============================================

// 是否给owner发过消息
const (
	PeerNotAck = 0 // 未发过
	PeerAck    = 1 // 发过
)
