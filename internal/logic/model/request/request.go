package request

type PingReq struct {
	Id int `json:"id"`
}

type SendMsgReq struct {
	SendId        uint64   `json:"send_id"`
	SendType      int32    `json:"send_type"`
	ReceiveId     uint64   `json:"receive_id"`
	ReceiveType   int32    `json:"receive_type"`
	MsgType       int32    `json:"msg_type"`
	Content       string   `json:"content"`
	InvisibleList []uint64 `json:"invisible_list"`
	SeqId         int64    `json:"seq_id"`
}

type FetchMsgReq struct {
	OwnerId   uint64 `json:"owner_id"`
	PeerId    uint64 `json:"peer_id"`
	PeerType  int32  `json:"peer_type"`
	FetchType string `json:"fetch_type"`
	VersionId uint64 `json:"version_id"`
}
