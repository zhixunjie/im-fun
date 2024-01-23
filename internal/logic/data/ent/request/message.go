package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
)

// MessageSendReq 发送消息给某个用户
type MessageSendReq struct {
	SeqId         model.BigIntType    `json:"seq_id"`         // 消息唯一id（客户端）
	SenderId      model.BigIntType    `json:"sender_id"`      // 消息发送者id
	SenderType    model.ContactIdType `json:"sender_type"`    // 消息发送者的用户类型
	ReceiverId    model.BigIntType    `json:"receiver_id"`    // 消息接收者id
	ReceiverType  model.ContactIdType `json:"receiver_type"`  // 消息接收者的用户类型
	MsgBody       format.MsgBody      `json:"msg_body"`       // 消息体
	InvisibleList []uint64            `json:"invisible_list"` // 不可见的列表
}

// MessageFetchReq 拉取消息列表（by version_id）
type MessageFetchReq struct {
	FetchType model.FetchType     `json:"fetch_type"`
	VersionId model.BigIntType    `json:"version_id"` // 版本id
	OwnerId   model.BigIntType    `json:"owner_id"`   // 会话拥有者
	OwnerType model.ContactIdType `json:"owner_type"` // 会话拥有者的用户类型
	PeerId    model.BigIntType    `json:"peer_id"`    // 会话联系人（对方）
	PeerType  model.ContactIdType `json:"peer_type"`  // 会话联系人（对方）的用户类型
}
