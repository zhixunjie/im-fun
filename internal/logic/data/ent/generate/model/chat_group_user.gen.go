// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameChatGroupUser = "chat_group_user"

// ChatGroupUser 群组与用户的绑定关系
type ChatGroupUser struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增id,主键" json:"id"`                   // 自增id,主键
	MemberID  uint64    `gorm:"column:member_id;not null;comment:群组成员id" json:"member_id"`                           // 群组成员id
	GroupID   uint64    `gorm:"column:group_id;not null;comment:群组id" json:"group_id"`                               // 群组id
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName ChatGroupUser's table name
func (*ChatGroupUser) TableName() string {
	return TableNameChatGroupUser
}