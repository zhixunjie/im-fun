package model

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/env"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
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

// TbNameMessage
// 因为 msgId 和 largerId 的后4位是相同的，所以这里传入 msgId 或者 largerId 都可以
func TbNameMessage(id uint64) (dbName string, tbName string) {
	// TODO: 临时测试
	return "im", TableNameChatMessage

	dbName = fmt.Sprintf("%v_%v", DbNameMessage, id%gen_id.SlotBit%DBNum())
	tbName = fmt.Sprintf("%v_%v", TableNameChatMessage, id%gen_id.SlotBit%TbNum())

	return dbName, tbName
}

func TbNameContact(ownerId uint64) (dbName string, tbName string) {
	// TODO: 临时测试
	return "im", TableNameChatContact

	dbName = fmt.Sprintf("%v_%v", DbNameMessage, ownerId%gen_id.SlotBit%DBNum())
	tbName = fmt.Sprintf("%v_%v", TableNameChatContact, ownerId%gen_id.SlotBit%TbNum())

	return dbName, tbName
}

func TbNameMessageByCId(id1, id2 *gmodel.ComponentId) (dbName string, tbName string) {
	switch {
	case id1.IsGroup(): // 群聊
		dbName, tbName = TbNameMessage(id1.Id())
	case id2.IsGroup(): // 群聊
		dbName, tbName = TbNameMessage(id2.Id())
	default: // 单聊
		_, largerId := id1.Sort(id2)
		dbName, tbName = TbNameMessage(largerId.Id())
	}

	return
}

// BigIntType 各种Id的类型（方便切换为int64、uint64）
type BigIntType = uint64

// =========================

// FetchMsgRangeParams 拉取消息列表
type FetchMsgRangeParams struct {
	FetchType                           gmodel.FetchType
	SessionId                           gmodel.SessionId
	LastDelMsgVersionId, PivotVersionId BigIntType // 确定消息的允许获取范围
	Limit                               int
	Owner, Peer                         *gmodel.ComponentId
}

// FetchContactRangeParams 拉取会话列表
type FetchContactRangeParams struct {
	FetchType      gmodel.FetchType
	Owner          *gmodel.ComponentId
	PivotVersionId BigIntType
	Limit          int
}

// BuildContactParams 构建参数
type BuildContactParams struct {
	Owner *gmodel.ComponentId
	Peer  *gmodel.ComponentId
}

// UpdateLastMsgId 更新最后一条消息
//type (
//	UpdateLastMsgId struct {
//		NewSessionId string
//		LastMsgId uint64
//		Peer1     UpdateLastMsgIdItem `json:"peer_1"`
//		Peer2     UpdateLastMsgIdItem `json:"peer_2"`
//	}
//	UpdateLastMsgIdItem struct {
//		ContactId uint64
//		OwnerId   *gen_id.ComponentId
//	}
//)
