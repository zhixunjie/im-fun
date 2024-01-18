package model

const (
	TotalDb           = 10
	TotalTableMessage = 100 // message表：分表个数（一共10个数据库，每个数据库100个表）
	TotalTableContact = 100 // contact表：分表个数（一共10个数据库，每个数据库100个表）
)

// 消息读取状态
const (
	MsgNotRead = 0 // 未读
	MsgRead    = 1 // 已读
)

// MsgStatus 消息状态
type MsgStatus uint32

const (
	MsgStatusNormal   MsgStatus = iota // 正常
	MsgStatusDel                       // 删除
	MsgStatusWithdraw                  // 后台删除
)

// FetchType 消息拉取方式
type FetchType = int32

const (
	FetchTypeBackward FetchType = iota // 拉取历史消息
	FetchTypeForward                   // 拉取最新消息
	FetchTypeBg                        // 后台拉消息（不清除未读数(history)）
)

type QueryMsgParams struct {
	FetchType           FetchType
	SmallerId           uint64
	LargerId            uint64
	PivotVersionId      uint64
	LastDelMsgVersionId uint64
	Limit               int
}

type QueryContactParams struct {
	FetchType      FetchType
	OwnerId        uint64
	PivotVersionId uint64
	Limit          int
}

type BuildContactParams struct {
	MsgId    uint64
	OwnerId  uint64
	PeerId   uint64
	PeerType PeerType
	PeerAck  uint32
}
