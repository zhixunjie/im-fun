package response

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
)

// MessageSendRsp 发送消息给某个用户
type MessageSendRsp struct {
	Base
	Data SendMsgRespData `json:"data"`
}

type SendMsgRespData struct {
	MsgId       uint64 `json:"msg_id"`
	SeqId       uint64 `json:"seq_id"`
	VersionId   uint64 `json:"version_id"`
	SortKey     uint64 `json:"sort_key"`
	SessionId   string `json:"session_id"`
	UnreadCount int64  `json:"unread_count"`
}

// MessageFetchRsp 拉取消息列表（by version_id）
type MessageFetchRsp struct {
	Base
	Data FetchMsgData `json:"data"`
}

type FetchMsgData struct {
	MsgList       []*MsgEntity `json:"msg_list"`        // 消息列表
	NextVersionId uint64       `json:"next_version_id"` // 最大的版本id
	HasMore       bool         `json:"has_more"`
}

type MsgEntity struct {
	MsgID     uint64              `json:"msg_id"`
	SeqID     uint64              `json:"seq_id"`
	MsgBody   format.MsgBody      `json:"msg_body"`
	SessionID string              `json:"session_id"`
	SenderID  uint64              `json:"sender_id"`
	VersionID uint64              `json:"version_id"`
	SortKey   uint64              `json:"sort_key"`
	Status    model.MsgStatus     `json:"status"`
	HasRead   model.MsgReadStatus `json:"has_read"`
}

// MessageWithdrawRsp 撤回消息
type MessageWithdrawRsp struct {
	Base
}

// DelBothSideRsp 删除消息（两边的聊天记录都需要删除）
type DelBothSideRsp struct {
	Base
}

// DelOneSideRsp 删除消息（只删除一边的聊天）
type DelOneSideRsp struct {
	Base
}

type ClearHistoryRsp struct {
	Base
}
