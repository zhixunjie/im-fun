package gen_id

import "fmt"

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
