package model

import "github.com/zhixunjie/im-fun/pkg/gen_id"

const (
	TotalDb           = 10
	TotalTableMessage = 512 // message表：分表个数（一共10个数据库，每个数据库512个表）
	TotalTableContact = 512 // contact表：分表个数（一共10个数据库，每个数据库512个表）
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
	MsgStatusWithdraw                  // 已撤回（双方都展示为撤回）
	MsgStatusDeleted                   // 已删除（双方都展示为删除）
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

// FetchMsgRangeParams 拉取消息列表
type FetchMsgRangeParams struct {
	FetchType                           FetchType
	SessionId                           string
	LastDelMsgVersionId, PivotVersionId BigIntType // 确定消息的允许获取范围
	Limit                               int
	OwnerId, PeerId                     *gen_id.ComponentId
}

// FetchContactRangeParams 拉取会话列表
type FetchContactRangeParams struct {
	FetchType      FetchType
	OwnerId        *gen_id.ComponentId
	PivotVersionId BigIntType
	Limit          int
}

// BuildContactParams 构建参数
type BuildContactParams struct {
	OwnerId *gen_id.ComponentId
	PeerId  *gen_id.ComponentId
}

// UpdateLastMsgId 更新最后一条消息
type UpdateLastMsgId struct {
	SessionId string
	LastMsgId uint64
	Peer1     UpdateLastMsgIdItem `json:"peer_1"`
	Peer2     UpdateLastMsgIdItem `json:"peer_2"`
}

type UpdateLastMsgIdItem struct {
	ContactId uint64
	OwnerId   *gen_id.ComponentId
}
