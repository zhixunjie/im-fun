package response

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
)

type PingResp struct {
	Base
	Pong string `json:"pong"`
}

// SendMsgResp 发送消息-响应
type SendMsgResp struct {
	Base
	Data SendMsgRespData `json:"data"`
}

type SendMsgRespData struct {
	MsgId       uint64 `json:"msg_id"`
	SeqId       uint64 `json:"seq_id"`
	VersionId   uint64 `json:"version_id"`
	SortKey     uint64 `json:"sort_key"`
	UnreadCount int64  `json:"unread_count"`
}

type FetchMsgResp struct {
	Base
	Data FetchMsgData `json:"data"`
}

type FetchMsgData struct {
	MsgList       []*Msg `json:"msg_list"`        // 获取得到的所有消息
	NextVersionId uint64 `json:"next_version_id"` // 最大的版本ID
	HasMore       bool   `json:"has_more"`
}

type Msg struct {
	MsgID     uint64         `json:"msg_id"`
	SeqID     uint64         `json:"seq_id"`
	MsgBody   format.MsgBody `json:"msg_body"`
	SessionID string         `json:"session_id"`
	SenderID  uint64         `json:"sender_id"`
	VersionID uint64         `json:"version_id"`
	SortKey   uint64         `json:"sort_key"`
	Status    uint32         `json:"status"`
	HasRead   uint32         `json:"has_read"`
}

type FetchSessionResp struct {
	Base
	Data FetchSessionData `json:"data"`
}

type FetchSessionData struct {
	ContactList   []*Contact `json:"contact_list"`
	NextVersionId uint64     `json:"next_version_id"` // 最大的版本ID
	HasMore       bool       `json:"has_more"`
}

type Contact struct {
	OwnerID      uint64 `json:"owner_id"`
	PeerID       uint64 `json:"peer_id"`
	PeerType     int32  `json:"peer_type"`
	PeerAck      uint32 `json:"peer_ack"`
	LastMsg      *Msg   `json:"last_msg"`
	VersionID    uint64 `json:"version_id"`
	SortKey      uint64 `json:"sort_key"`
	Status       uint32 `json:"status"`
	Labels       string `json:"labels"`
	UnreadMsgNum int64  `json:"unread_msg_num"` // 当前会话框的未读信息数
}
