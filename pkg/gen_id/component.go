package gen_id

import (
	"fmt"
	"github.com/samber/lo"
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

func NewComponentId(id uint64, idType ContactIdType) *ComponentId {
	return &ComponentId{
		id:     id,
		idType: idType,
	}
}

// Sort 小的id在前，大的id在后
func Sort(a, b *ComponentId) (*ComponentId, *ComponentId) {
	if a.id < b.id {
		return a, b
	}

	return b, a
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
