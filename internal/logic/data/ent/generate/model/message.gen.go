// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMessage = "message"

// Message 消息表（一条私信只有一行记录）
type Message struct {
	MsgID         uint64    `gorm:"column:msg_id;primaryKey;comment:消息唯一id（服务端生成）" json:"msg_id"`                             // 消息唯一id（服务端生成）
	SeqID         uint64    `gorm:"column:seq_id;not null;comment:消息唯一id（客户端生成）" json:"seq_id"`                               // 消息唯一id（客户端生成）
	MsgType       uint32    `gorm:"column:msg_type;not null;comment:消息类型" json:"msg_type"`                                    // 消息类型
	Content       string    `gorm:"column:content;not null;comment:消息内容，json格式" json:"content"`                               // 消息内容，json格式
	SessionID     string    `gorm:"column:session_id;not null;comment:会话id" json:"session_id"`                                // 会话id
	SenderID      uint64    `gorm:"column:sender_id;not null;comment:私信发送者id" json:"sender_id"`                               // 私信发送者id
	VersionID     uint64    `gorm:"column:version_id;not null;comment:版本id（用于拉取消息）" json:"version_id"`                        // 版本id（用于拉取消息）
	SortKey       uint64    `gorm:"column:sort_key;not null;comment:消息展示顺序（按顺序展示消息）" json:"sort_key"`                         // 消息展示顺序（按顺序展示消息）
	Status        uint32    `gorm:"column:status;not null;comment:消息状态。0：正常，1：已删除，2：已撤回" json:"status"`                       // 消息状态。0：正常，1：已删除，2：已撤回
	HasRead       uint32    `gorm:"column:has_read;not null;comment:接收方是否已读，0：未读，1：已读" json:"has_read"`                       // 接收方是否已读，0：未读，1：已读
	InvisibleList string    `gorm:"column:invisible_list;not null;comment:发送方看到消息发出去了，但是对于在列表的用户是不可见的" json:"invisible_list"` // 发送方看到消息发出去了，但是对于在列表的用户是不可见的
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`      // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`      // 更新时间
}

// TableName Message's table name
func (*Message) TableName() string {
	return TableNameMessage
}
