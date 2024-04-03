package response

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

// ContactFetchRsp 拉取会话列表（by version_id）
type ContactFetchRsp struct {
	Base
	Data FetchContactData `json:"data"`
}

type FetchContactData struct {
	ContactList   []*ContactEntity `json:"contact_list"`    // 联系人列表
	NextVersionId model.BigIntType `json:"next_version_id"` // 最大的版本ID
	HasMore       bool             `json:"has_more"`
}

type ContactEntity struct {
	OwnerID      model.BigIntType     `json:"owner_id"`
	OwnerType    gen_id.ContactIdType `json:"owner_type"`
	PeerID       model.BigIntType     `json:"peer_id"`
	PeerType     gen_id.ContactIdType `json:"peer_type"`
	PeerAck      model.PeerAckStatus  `json:"peer_ack"`
	VersionID    model.BigIntType     `json:"version_id"`
	SortKey      model.BigIntType     `json:"sort_key"`
	Status       model.ContactStatus  `json:"status"`
	Labels       string               `json:"labels"`
	LastMsg      *MsgEntity           `json:"last_msg"`
	UnreadMsgNum int64                `json:"unread_msg_num"` // 当前会话框的未读信息数
}
