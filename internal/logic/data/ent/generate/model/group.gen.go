// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameGroup = "group"

// Group 群组表
type Group struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增id,主键" json:"id"`                   // 自增id,主键
	Name      string    `gorm:"column:name;not null;comment:群组名称" json:"name"`                                       // 群组名称
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Group's table name
func (*Group) TableName() string {
	return TableNameGroup
}
