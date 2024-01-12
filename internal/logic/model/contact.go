package model

import "time"

// 联系人状态
const (
	ContactStatusNormal = 0 // 正常
	ContactStatusDel    = 1 // 删除
)

// 联系人类型
// 0-99业务自己扩展，100之后保留
const (
	PeerNotExist = -1
	PeerNormal   = 0   // 普通用户（peer_id等于用户id）
	PeerSys      = 100 // 系统用户（peer_id等于100000）
	PeerGroup    = 101 // 群组（peer_id等于群组id）
)

// 是否给owner发过消息
const (
	PeerNotAck = 0 // 未发过
	PeerAck    = 1 // 发过
)

type Contact struct {
	Id           uint64    `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"` // 自增id（主键）
	OwnerId      uint64    `json:"owner_id" gorm:"column:owner_id"`                // 会话拥有者
	PeerId       uint64    `json:"peer_id" gorm:"column:peer_id"`                  // 联系人（对方用户）
	PeerType     int32     `json:"peer_type" gorm:"column:peer_type"`              // 联系人类型
	PeerAck      int32     `json:"peer_ack" gorm:"column:peer_ack"`                // peer是否给owner发过消息，0：未发过，1：发过
	LastMsgId    uint64    `json:"last_msg_id" gorm:"column:last_msg_id"`          // 聊天记录中，最新一条发送的私信id
	LastDelMsgId uint64    `json:"last_del_msg_id" gorm:"column:last_del_msg_id"`  // 聊天记录中，最后一次删除联系人时的私信id
	VersionId    uint64    `json:"version_id" gorm:"column:version_id"`            // 版本ID（用于拉取会话框）
	SortKey      uint64    `json:"sort_key" gorm:"column:sort_key"`                // 会话展示顺序（按顺序展示会话），可修改顺序，如：置顶操作
	Status       int32     `json:"status" gorm:"column:status"`                    // 联系人状态，0：正常，1：被删除
	Labels       string    `json:"labels" gorm:"column:labels"`                    // 会话标签，json格式
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
