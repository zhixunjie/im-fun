package response

import "github.com/zhixunjie/im-fun/internal/logic/data/ent/format"

// MessageSendRsp 发送消息-响应
type MessageSendRsp struct {
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

type MessageFetchRsp struct {
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
