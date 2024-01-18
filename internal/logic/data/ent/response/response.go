package response

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
	Data string
}
