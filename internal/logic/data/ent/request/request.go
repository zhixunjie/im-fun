package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
)

type PingReq struct {
	Id int `json:"id"`
}

// MessageSendReq 发送消息-请求
type MessageSendReq struct {
	SeqId                   uint64                `json:"seq_id"`                     // 消息唯一ID（客户端）
	SenderId                uint64                `json:"sender_id"`                  // 消息发送者ID
	ReceiverId              uint64                `json:"receiver_id"`                // 消息接收者ID
	SenderContactPeerType   model.ContactPeerType `json:"sender_contact_peer_type"`   // 消息发送者的联系人类型
	ReceiverContactPeerType model.ContactPeerType `json:"receiver_contact_peer_type"` // 消息接收者的联系人类型
	InvisibleList           []uint64              `json:"invisible_list"`             // 不可见的列表
	MsgBody                 format.MsgBody        `json:"msg_body"`                   // 消息体
}

type MessageFetchReq struct {
	FetchType model.FetchType `json:"fetch_type"`
	OwnerId   uint64          `json:"owner_id"`
	PeerId    uint64          `json:"peer_id"`
	VersionId uint64          `json:"version_id"`
}

type ContactFetchReq struct {
	OwnerId   uint64 `json:"owner_id"`
	VersionId uint64 `json:"version_id"`
}
