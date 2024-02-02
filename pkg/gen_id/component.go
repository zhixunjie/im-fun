package gen_id

import (
	"fmt"
	"github.com/samber/lo"
)

// ComponentId 组合ID
type ComponentId struct {
	id     uint64
	idType uint32
}

func (c *ComponentId) ToString() string {
	return fmt.Sprintf("%d_%d", c.idType, c.id)
}

func (c *ComponentId) Id() uint64 {
	return c.id
}

func (c *ComponentId) Type() uint32 {
	return c.idType
}

func (c *ComponentId) IsGroup() bool {
	typeArr := []uint32{uint32(ContactIdTypeGroup)}

	return lo.Contains(typeArr, c.idType)
}

func (c *ComponentId) Equal(b *ComponentId) bool {
	if c.Id() == b.Id() && c.Type() == b.Type() {
		return true
	}
	return false
}

func NewComponentId(id uint64, idType uint32) *ComponentId {
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
	return NewComponentId(id, uint32(ContactIdTypeUser))
}

// NewRobotComponentId 用户ID
func NewRobotComponentId(id uint64) *ComponentId {
	return NewComponentId(id, uint32(ContactIdTypeRobot))
}

// NewSystemComponentId 用户ID
func NewSystemComponentId(id uint64) *ComponentId {
	return NewComponentId(id, uint32(ContactIdTypeSystem))
}

// NewGroupComponentId 用户ID
func NewGroupComponentId(id uint64) *ComponentId {
	return NewComponentId(id, uint32(ContactIdTypeGroup))
}
