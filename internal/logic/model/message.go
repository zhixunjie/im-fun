package model

import "time"

const (
	TotalDb           = 10
	TotalTableMessage = 100 // message表：分表个数
	TotalTableContact = 100 // contact表：分表个数
)

// 消息读取状态
const (
	MsgNotRead = 0 // 未读
	MsgRead    = 1 // 已读
)

// 消息状态
const (
	MsgStatusNormal   = 0 // 正常
	MsgStatusDel      = 1 // 删除
	MsgStatusWithdraw = 2 // 后台删除
)

type Message struct {
	MsgId         uint64    `json:"msg_id" gorm:"PRIMARY_KEY;column:msg_id"`     // 消息id
	MsgType       int32     `json:"msg_type" gorm:"column:msg_type"`             // 消息类型
	SessionId     string    `json:"session_id" gorm:"column:session_id"`         // 会话id
	SendId        uint64    `json:"send_id" gorm:"column:send_id"`               // 消息发送者
	VersionId     uint64    `json:"version_id" gorm:"column:version_id"`         // 版本号（用于拉取消息）
	SortKey       uint64    `json:"sort_key" gorm:"column:sort_key"`             // 消息展示顺序（按顺序展示消息）
	Status        int32     `json:"status" gorm:"column:status"`                 // 消息状态。0：正常，1：被审核删除，2：撤销
	Content       string    `json:"content" gorm:"column:content"`               // 消息内容
	HasRead       int32     `json:"has_read" gorm:"column:has_read"`             // 消息接收者是否已读消息。0：未读，1：已读
	InvisibleList string    `json:"invisible_list" gorm:"column:invisible_list"` // 消息发出去了，但是对于在列表的用户是不可见的
	SeqId         int64     `json:"seq_id" gorm:"column:seq_id"`                 // 客户端本地数据库的消息id
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
}
