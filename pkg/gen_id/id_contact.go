package gen_id

import (
	"fmt"
	"github.com/samber/lo"
)

// note: 这个文件可以放全局管理

// ContactIdType 联系人类型
// 1-99业务自己扩展，100之后保留
type ContactIdType uint32

const (
	TypeUser   ContactIdType = 1   // 对方是普通用户
	TypeRobot  ContactIdType = 2   // 对方是机器人
	TypeSystem ContactIdType = 100 // 对方是系统用户
	TypeGroup  ContactIdType = 101 // 对方是群组
)

// ComponentId 组合ID
type ComponentId struct {
	id     uint64
	idType ContactIdType
}

func (c *ComponentId) ToString() string {
	return fmt.Sprintf("%d_%d", c.idType, c.id)
}

func (c *ComponentId) Id() uint64 {
	return c.id
}

func (c *ComponentId) Type() ContactIdType {
	return c.idType
}

func (c *ComponentId) IsGroup() bool {
	typeArr := []ContactIdType{TypeGroup}

	return lo.Contains(typeArr, c.idType)
}

func (c *ComponentId) Equal(b *ComponentId) bool {
	if c.Id() == b.Id() && c.Type() == b.Type() {
		return true
	}
	return false
}

// Sort 小的id在前，大的id在后
func (c *ComponentId) Sort(b *ComponentId) (*ComponentId, *ComponentId) {
	if c.id < b.id {
		return c, b
	}

	return b, c
}

func NewComponentId(id uint64, idType ContactIdType) *ComponentId {
	return &ComponentId{
		id:     id,
		idType: idType,
	}
}

// 预定义的组合ID

// NewUserComponentId 用户ID
func NewUserComponentId(id uint64) *ComponentId {
	return NewComponentId(id, TypeUser)
}

// NewRobotComponentId 机器人ID
func NewRobotComponentId(id uint64) *ComponentId {
	return NewComponentId(id, TypeRobot)
}

// NewSystemComponentId 系统ID
func NewSystemComponentId(id uint64) *ComponentId {
	return NewComponentId(id, TypeSystem)
}

// NewGroupComponentId 群组ID
func NewGroupComponentId(id uint64) *ComponentId {
	return NewComponentId(id, TypeGroup)
}
