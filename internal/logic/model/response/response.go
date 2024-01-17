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
	MsgId        uint64 `json:"msg_id"`
	SeqId        int64  `json:"seq_id"`
	MsgVersionId uint64 `json:"msg_version_id"`
	MsgSortKey   uint64 `json:"msg_sort_key"`
	UnreadCount  int64  `json:"unread_count"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time"`
}

type FetchMsgResp struct {
	Base
	Data string
}
