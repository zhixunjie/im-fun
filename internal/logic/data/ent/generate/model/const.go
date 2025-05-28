package model

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/env"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

const (
	DbNameMessage = "im_message"
)

func DBNum() (num uint64) {
	num = 10
	if !env.IsProd() {
		num = 2
		return
	}
	return
}

func TbNum() (num uint64) {
	num = 512
	if !env.IsProd() {
		num = 4
		return
	}
	return
}

// ShardingTbNameMessage
// 因为 msgId 和 largerId 的后4位是相同的，所以这里传入 msgId 或者 largerId 都可以
func ShardingTbNameMessage(id uint64) (dbName string, tbName string) {
	dbName = fmt.Sprintf("%v_%v", DbNameMessage, id%gen_id.SlotBit%DBNum())
	tbName = fmt.Sprintf("%v_%v", TableNameMessage, id%gen_id.SlotBit%TbNum())

	return dbName, tbName
}

func ShardingTbNameMessageByComponentId(id1, id2 *gen_id.ComponentId) (dbName string, tbName string) {
	switch {
	case id1.IsGroup(): // 群聊
		dbName, tbName = ShardingTbNameMessage(id1.Id())
	case id2.IsGroup(): // 群聊
		dbName, tbName = ShardingTbNameMessage(id2.Id())
	default: // 单聊
		_, largerId := id1.Sort(id2)
		dbName, tbName = ShardingTbNameMessage(largerId.Id())
	}

	return
}

func ShardingTbNameContact(ownerId uint64) (dbName string, tbName string) {
	dbName = fmt.Sprintf("%v_%v", DbNameMessage, ownerId%gen_id.SlotBit%DBNum())
	tbName = fmt.Sprintf("%v_%v", TableNameContact, ownerId%gen_id.SlotBit%TbNum())

	return dbName, tbName
}

// BigIntType 各种Id的类型（方便切换为int64、uint64）
type BigIntType = uint64

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

// =========================

// FetchMsgRangeParams 拉取消息列表
type FetchMsgRangeParams struct {
	FetchType                           FetchType
	SessionId                           string
	LastDelMsgVersionId, PivotVersionId BigIntType // 确定消息的允许获取范围
	Limit                               int
	Owner, Peer                         *gen_id.ComponentId
}

// FetchContactRangeParams 拉取会话列表
type FetchContactRangeParams struct {
	FetchType      FetchType
	Owner          *gen_id.ComponentId
	PivotVersionId BigIntType
	Limit          int
}

// BuildContactParams 构建参数
type BuildContactParams struct {
	Owner *gen_id.ComponentId
	Peer  *gen_id.ComponentId
}

// UpdateLastMsgId 更新最后一条消息
//type (
//	UpdateLastMsgId struct {
//		SessionId string
//		LastMsgId uint64
//		Peer1     UpdateLastMsgIdItem `json:"peer_1"`
//		Peer2     UpdateLastMsgIdItem `json:"peer_2"`
//	}
//	UpdateLastMsgIdItem struct {
//		ContactId uint64
//		OwnerId   *gen_id.ComponentId
//	}
//)
