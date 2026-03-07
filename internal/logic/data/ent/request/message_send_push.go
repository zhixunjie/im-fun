package request

import "github.com/zhixunjie/im-fun/internal/logic/data/ent/format"

// MessageSendWithPushReq 发送消息并推送长链接（鉴权接口，发送方从 token 中解析）
type MessageSendWithPushReq struct {
	ReceiverUniId string          `json:"receiver_uni_id"` // 接收方用户ID
	MsgBody       *format.MsgBody `json:"msg_body"`        // 消息体
}
