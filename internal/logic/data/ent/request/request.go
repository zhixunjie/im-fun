package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
)

type PingReq struct {
	Id int `json:"id"`
}

// MessageSendReq 发送消息给某个用户
type MessageSendReq struct {
	SeqId                   model.BigIntType      `json:"seq_id"`                     // 消息唯一ID（客户端）
	SenderId                model.BigIntType      `json:"sender_id"`                  // 消息发送者ID
	ReceiverId              model.BigIntType      `json:"receiver_id"`                // 消息接收者ID
	SenderContactPeerType   model.ContactPeerType `json:"sender_contact_peer_type"`   // 消息发送者的联系人类型
	ReceiverContactPeerType model.ContactPeerType `json:"receiver_contact_peer_type"` // 消息接收者的联系人类型
	MsgBody                 format.MsgBody        `json:"msg_body"`                   // 消息体
	InvisibleList           []uint64              `json:"invisible_list"`             // 不可见的列表
}

// MessageFetchReq 拉取消息列表（by version_id）
type MessageFetchReq struct {
	FetchType model.FetchType       `json:"fetch_type"`
	VersionId model.BigIntType      `json:"version_id"` // 版本ID
	OwnerId   model.BigIntType      `json:"owner_id"`   // 会话拥有者
	PeerId    model.BigIntType      `json:"peer_id"`    // 会话联系人
	PeerType  model.ContactPeerType `json:"peer_type"`  // 会话联系人的联系人类型
}

// ContactFetchReq 拉取会话列表（by version_id）
type ContactFetchReq struct {
	VersionId model.BigIntType `json:"version_id"` // 版本ID
	OwnerId   model.BigIntType `json:"owner_id"`   // 会话拥有者
}
