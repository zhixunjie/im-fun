package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

// MessageSendReq 发送消息给某个用户
type MessageSendReq struct {
	SeqId         model.BigIntType     `json:"seq_id"`         // 消息唯一id（客户端）
	SenderId      model.BigIntType     `json:"sender_id"`      // 消息发送者id
	SenderType    gen_id.ContactIdType `json:"sender_type"`    // 消息发送者的用户类型
	ReceiverId    model.BigIntType     `json:"receiver_id"`    // 消息接收者id
	ReceiverType  gen_id.ContactIdType `json:"receiver_type"`  // 消息接收者的用户类型
	MsgBody       format.MsgBody       `json:"msg_body"`       // 消息体
	InvisibleList []uint64             `json:"invisible_list"` // 不可见的列表
}

// MessageFetchReq 拉取消息列表（by version_id）
type MessageFetchReq struct {
	FetchType model.FetchType      `json:"fetch_type"`
	VersionId model.BigIntType     `json:"version_id"` // 版本id
	OwnerId   model.BigIntType     `json:"owner_id"`   // 会话拥有者
	OwnerType gen_id.ContactIdType `json:"owner_type"` // 会话拥有者的用户类型
	PeerId    model.BigIntType     `json:"peer_id"`    // 会话联系人（对方）
	PeerType  gen_id.ContactIdType `json:"peer_type"`  // 会话联系人（对方）的用户类型
}

// MessageWithdrawReq 撤回消息
type MessageWithdrawReq struct {
	MsgId      model.BigIntType     `json:"msg_id"`      // 撤回哪一条消息？
	SenderId   model.BigIntType     `json:"sender_id"`   // 消息发送者id
	SenderType gen_id.ContactIdType `json:"sender_type"` // 消息发送者的用户类型
}

// DelBothSideReq 删除消息（两边的聊天记录都需要删除）
type DelBothSideReq struct {
	MsgId      model.BigIntType     `json:"msg_id"`      // 删除哪一条消息？
	SenderId   model.BigIntType     `json:"sender_id"`   // 消息发送者id
	SenderType gen_id.ContactIdType `json:"sender_type"` // 消息发送者的用户类型
}

// DelOneSideReq 删除消息（只删除一边的聊天）
type DelOneSideReq struct {
	MsgId      model.BigIntType     `json:"msg_id"`      // 删除哪一条消息？
	SenderId   model.BigIntType     `json:"sender_id"`   // 消息发送者id
	SenderType gen_id.ContactIdType `json:"sender_type"` // 消息发送者的用户类型
}

// ClearHistoryReq 清空聊天记录
type ClearHistoryReq struct {
	MsgId     model.BigIntType     `json:"msg_id"`     // 从哪一条消息开始，进行聊天记录的清空
	OwnerId   model.BigIntType     `json:"owner_id"`   // 会话拥有者
	OwnerType gen_id.ContactIdType `json:"owner_type"` // 会话拥有者的用户类型
	PeerId    model.BigIntType     `json:"peer_id"`    // 会话联系人（对方）
	PeerType  gen_id.ContactIdType `json:"peer_type"`  // 会话联系人（对方）的用户类型
}
