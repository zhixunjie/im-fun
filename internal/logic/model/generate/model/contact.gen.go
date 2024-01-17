// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameContact = "contact"

// Contact 会话表（通信双方各有一行记录）
type Contact struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增id,主键" json:"id"`                    // 自增id,主键
	OwnerID      int64     `gorm:"column:owner_id;not null;comment:会话拥有者" json:"owner_id"`                               // 会话拥有者
	PeerID       int64     `gorm:"column:peer_id;not null;comment:联系人（对方用户）" json:"peer_id"`                             // 联系人（对方用户）
	PeerType     int32     `gorm:"column:peer_type;not null;comment:联系人类型，0：普通，100：系统，101：群组" json:"peer_type"`          // 联系人类型，0：普通，100：系统，101：群组
	PeerAck      int32     `gorm:"column:peer_ack;not null;comment:peer_id是否给owner发过消息，0：未发过，1：发过" json:"peer_ack"`      // peer_id是否给owner发过消息，0：未发过，1：发过
	LastMsgID    int64     `gorm:"column:last_msg_id;not null;comment:聊天记录中，最新一条发送的私信id" json:"last_msg_id"`             // 聊天记录中，最新一条发送的私信id
	LastDelMsgID int64     `gorm:"column:last_del_msg_id;not null;comment:聊天记录中，最后一次删除联系人时的私信id" json:"last_del_msg_id"` // 聊天记录中，最后一次删除联系人时的私信id
	VersionID    int64     `gorm:"column:version_id;not null;comment:版本id（用于拉取会话框）" json:"version_id"`                   // 版本id（用于拉取会话框）
	SortKey      int64     `gorm:"column:sort_key;not null;comment:会话展示顺序（按顺序展示会话）可修改顺序，如：置顶操作" json:"sort_key"`         // 会话展示顺序（按顺序展示会话）可修改顺序，如：置顶操作
	Status       int32     `gorm:"column:status;not null;comment:联系人状态，0：正常，1：被删除" json:"status"`                        // 联系人状态，0：正常，1：被删除
	Labels       string    `gorm:"column:labels;not null;comment:会话标签，json字符串" json:"labels"`                            // 会话标签，json字符串
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`  // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`  // 更新时间
}

// TableName Contact's table name
func (*Contact) TableName() string {
	return TableNameContact
}
