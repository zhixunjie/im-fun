package response

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"time"
)

// MessageSendRsp 发送消息给某个用户
type MessageSendRsp struct {
	Base
	Data *SendMsgRespData `json:"data"`
}

type SendMsgRespData struct {
	MsgID       uint64 `json:"msg_id,string"`
	SeqID       uint64 `json:"seq_id,string"`
	VersionID   uint64 `json:"version_id,string"`
	SortKey     uint64 `json:"sort_key,string"`
	SessionId   string `json:"session_id"`
	UnreadCount int64  `json:"unread_count"`
}

// MessageFetchRsp 拉取消息列表（by version_id）
type MessageFetchRsp struct {
	Base
	Data *FetchMsgData `json:"data"`
}

type FetchMsgData struct {
	HasMore       bool         `json:"has_more"`
	NextVersionId uint64       `json:"next_version_id,string"` // 最大的版本id
	MsgList       []*MsgEntity `json:"msg_list"`               // 消息列表
}

type MsgEntity struct {
	MsgID     uint64               `json:"msg_id,string"`
	SeqID     uint64               `json:"seq_id,string"`
	MsgBody   *format.MsgBody      `json:"msg_body"`
	SessionID string               `json:"session_id"`
	SenderID  uint64               `json:"sender_id"`
	SendType  gmodel.ContactIdType `json:"send_type"`
	VersionID uint64               `json:"version_id,string"`
	SortKey   uint64               `json:"sort_key,string"`
	Status    gmodel.MsgStatus     `json:"status"`
	HasRead   gmodel.MsgReadStatus `json:"has_read"`
	CreatedAt time.Time            `json:"created_at"` // 创建时间
	UpdatedAt time.Time            `json:"updated_at"` // 更新时间
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

type MessageClearHistoryRsp struct {
	Base
	LastDelMsgID uint64 `json:"last_del_msg_id,string"` // 最后一条删除的消息id
}
