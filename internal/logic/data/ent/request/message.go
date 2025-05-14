package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

// MessageSendReq 发送消息给某个用户
type MessageSendReq struct {
	SeqId         model.BigIntType    `json:"seq_id,string"`  // 消息唯一id（客户端）
	Sender        *gen_id.ComponentId `json:"sender"`         // 消息发送者
	Receiver      *gen_id.ComponentId `json:"receiver"`       // 消息接收者
	MsgBody       format.MsgBody      `json:"msg_body"`       // 消息体
	InvisibleList []uint64            `json:"invisible_list"` // 不可见的列表
}

// MessageFetchReq 拉取消息列表（by version_id）
type MessageFetchReq struct {
	FetchType model.FetchType     `json:"fetch_type"`
	VersionId model.BigIntType    `json:"version_id,string"` // 版本id
	Owner     *gen_id.ComponentId `json:"owner"`             // 会话拥有者
	Peer      *gen_id.ComponentId `json:"peer"`              // 会话联系人（对方）
}

// MessageRecallReq 撤回消息
type MessageRecallReq struct {
	MsgID  model.BigIntType    `json:"msg_id,string"` // 撤回哪一条消息？
	Sender *gen_id.ComponentId `json:"sender"`        // 消息发送者
}

// MessageDelBothSideReq 删除消息（两边的聊天记录都需要删除）
type MessageDelBothSideReq struct {
	MsgID  model.BigIntType    `json:"msg_id,string"` // 删除哪一条消息？
	Sender *gen_id.ComponentId `json:"sender"`        // 消息发送者
}

// MessageDelOneSideReq 删除消息（只删除一边的聊天）
type MessageDelOneSideReq struct {
	MsgID  model.BigIntType    `json:"msg_id,string"` // 删除哪一条消息？
	Sender *gen_id.ComponentId `json:"sender"`        // 消息发送者
}

// ClearHistoryReq 清空聊天记录
type ClearHistoryReq struct {
	MsgID model.BigIntType    `json:"msg_id,string"` // 从哪一条消息开始，进行聊天记录的清空
	Owner *gen_id.ComponentId `json:"owner"`         // 会话拥有者
	Peer  *gen_id.ComponentId `json:"peer"`          // 会话联系人（对方）
}
