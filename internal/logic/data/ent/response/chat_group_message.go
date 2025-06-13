package response

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"time"
)

type (
	// GroupMessageSendRsp 发送消息给某个用户
	GroupMessageSendRsp struct {
		Base
		Data *GroupMessageSendData `json:"data"`
	}
	GroupMessageSendData struct {
		MsgID       uint64 `json:"msg_id,string"`
		SeqID       uint64 `json:"seq_id,string"`
		VersionID   uint64 `json:"version_id,string"`
		SortKey     uint64 `json:"sort_key,string"`
		SessionId   string `json:"session_id"`
		UnreadCount int64  `json:"unread_count"`
	}
)

type (
	// GroupMessageFetchRsp 拉取消息列表（by version_id）
	GroupMessageFetchRsp struct {
		Base
		Data *GroupMessageFetchData `json:"data"`
	}
	GroupMessageFetchData struct {
		HasMore       bool              `json:"has_more"`
		NextVersionId uint64            `json:"next_version_id,string"` // 最大的版本id
		MsgList       []*GroupMsgEntity `json:"msg_list"`               // 消息列表
	}
	// GroupMsgEntity 消息实体
	GroupMsgEntity struct {
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
)
