package model

const (
	TotalDb           = 10
	TotalTableMessage = 100 // message表：分表个数（一共10个数据库，每个数据库100个表）
	TotalTableContact = 100 // contact表：分表个数（一共10个数据库，每个数据库100个表）
)

// BigIntType 各种Id的类型（方便切换为int64、uint64）
type BigIntType = uint64

// ================================ Contact ================================

// ContactStatus 联系人状态
type ContactStatus uint32

const (
	ContactStatusNormal  = 0 // 正常
	ContactStatusDeleted = 1 // 已删除
)

// =========================

// ContactPeerType 联系人类型
// 0-99业务自己扩展，100之后保留
type ContactPeerType int32

const (
	PeerTypeNormalUser ContactPeerType = 0   // 会话框的对方是普通用户（peer_id等于用户id）
	PeerTypeSystemUser ContactPeerType = 100 // 会话框的对方是系统用户（peer_id等于 SystemUid ）
	PeerTypeGroup      ContactPeerType = 101 // 会话框的对方是群组（peer_id等于群组id）

	// SystemUid 系统用户的用户ID
	SystemUid = 100000
)

// =========================

// PeerAckStatus 是否给owner发过消息
type PeerAckStatus uint32

const (
	PeerNotAck PeerAckStatus = 0 // 未发过
	PeerAcked                = 1 // 发过
)

// ================================ Message ================================

// =========================

// MsgReadStatus 消息读取状态
type MsgReadStatus uint32

const (
	MsgNotRead MsgReadStatus = 0 // 未读
	MsgRead                  = 1 // 已读
)

// =========================

// MsgStatus 消息状态
type MsgStatus uint32

const (
	MsgStatusNormal   MsgStatus = iota // 正常
	MsgStatusDeleted                   // 已删除
	MsgStatusWithdraw                  // 已撤回
)

// =========================

// FetchType 消息拉取方式
type FetchType = int32

const (
	FetchTypeBackward FetchType = iota // 拉取历史消息
	FetchTypeForward                   // 拉取最新消息
	FetchTypeInBg                      // 后台拉消息（不清除未读数(history)）
)

// =========================

type FetchMsgRangeParams struct {
	FetchType           FetchType
	SmallerId           BigIntType
	LargerId            BigIntType
	PivotVersionId      BigIntType
	LastDelMsgVersionId BigIntType
	Limit               int
}

type FetchContactRangeParams struct {
	FetchType      FetchType
	OwnerId        BigIntType
	PivotVersionId BigIntType
	Limit          int
}

type BuildContactParams struct {
	OwnerId      BigIntType
	PeerId       BigIntType
	LastMsgId    BigIntType
	InitPeerType ContactPeerType
	InitPeerAck  uint32
}
