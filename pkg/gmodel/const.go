package gmodel

// ================================ Contact ================================

// ContactStatus 联系人状态
type ContactStatus uint32

const (
	ContactStatusNormal  ContactStatus = 1 // 正常
	ContactStatusDeleted ContactStatus = 2 // 已删除
)

// =========================

// PeerAckStatus 是否给owner发过消息
type PeerAckStatus uint32

const (
	PeerNotAck PeerAckStatus = 1 // 未回答过消息
	PeerAcked  PeerAckStatus = 2 // 回答过消息
)

// ================================ Message ================================

// MsgStatus 消息状态
type MsgStatus uint32

const (
	MsgStatusNormal  MsgStatus = 1 // 正常
	MsgStatusDeleted MsgStatus = 2 // 已删除（双方都展示为删除）
	MsgStatusRecall  MsgStatus = 3 // 已撤回（双方都展示为撤回）
)

// MsgReadStatus 消息读取状态
type MsgReadStatus uint32

const (
	MsgNotRead MsgReadStatus = 1 // 未读
	MsgRead    MsgReadStatus = 2 // 已读
)

// FetchType 消息拉取方式
type FetchType = int32

const (
	FetchTypeBackward FetchType = 1 // 拉取历史消息
	FetchTypeForward  FetchType = 2 // 拉取最新消息
	FetchTypeInBg     FetchType = 3 // 后台拉消息：不清除未读数(history)
)
