package response

import "github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"

type ContactFetchRsp struct {
	Base
	Data FetchContactData `json:"data"`
}

type FetchContactData struct {
	ContactList   []*Contact `json:"contact_list"`
	NextVersionId uint64     `json:"next_version_id"` // 最大的版本ID
	HasMore       bool       `json:"has_more"`
}

type Contact struct {
	OwnerID      uint64                `json:"owner_id"`
	PeerID       uint64                `json:"peer_id"`
	PeerType     model.ContactPeerType `json:"peer_type"`
	PeerAck      model.PeerAckStatus   `json:"peer_ack"`
	VersionID    uint64                `json:"version_id"`
	SortKey      uint64                `json:"sort_key"`
	Status       model.ContactStatus   `json:"status"`
	Labels       string                `json:"labels"`
	LastMsg      *Msg                  `json:"last_msg"`
	UnreadMsgNum int64                 `json:"unread_msg_num"` // 当前会话框的未读信息数
}
