package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
)

// GroupMessageSendReq 发送消息给某个用户
type GroupMessageSendReq struct {
	SeqId         model.BigIntType    `json:"seq_id,string"`  // 消息唯一id（客户端）
	Sender        *gmodel.ComponentId `json:"sender"`         // 消息发送者
	Receiver      *gmodel.ComponentId `json:"receiver"`       // 消息接收者
	MsgBody       *format.MsgBody     `json:"msg_body"`       // 消息体
	InvisibleList []string            `json:"invisible_list"` // 不可见的列表
}

// GroupMessageFetchReq 拉取消息列表（by version_id）
type GroupMessageFetchReq struct {
	FetchType gmodel.FetchType    `json:"fetch_type"`
	VersionId model.BigIntType    `json:"version_id,string"` // 版本id
	Owner     *gmodel.ComponentId `json:"owner"`             // 会话拥有者
	Peer      *gmodel.ComponentId `json:"peer"`              // 会话联系人（对方）
}
