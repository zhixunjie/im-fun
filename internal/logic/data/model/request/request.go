package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/model/format"
)

type PingReq struct {
	Id int `json:"id"`
}

// SendMsgReq 发送消息-请求
type SendMsgReq struct {
	SendId        uint64         `json:"send_id"`        // 发送者ID
	SendType      int32          `json:"send_type"`      // 发送者类型
	PeerId        uint64         `json:"peer_id"`        // 接收者ID
	PeerType      int32          `json:"peer_type"`      // 接收者类型
	InvisibleList []uint64       `json:"invisible_list"` // 不可见的列表
	SeqId         int64          `json:"seq_id"`         // 客户端序列ID
	MsgBody       format.MsgBody `json:"msg_body"`       // 消息体
}

type FetchMsgReq struct {
	OwnerId   uint64 `json:"owner_id"`
	PeerId    uint64 `json:"peer_id"`
	PeerType  int32  `json:"peer_type"`
	FetchType string `json:"fetch_type"`
	VersionId uint64 `json:"version_id"`
}
