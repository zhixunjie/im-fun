package bytes

import (
	"sync"
)

// PoolSync：使用sync包编写的BufferPool
// 优点：代码更加简洁
// 缺点：分配方式不够高效，发现Buffer不足时，只会一个个去New

// PoolSync
// A BufferPool based on sync.Pool
type PoolSync struct {
	pool sync.Pool
}

func (p *PoolSync) Init(bufNum, bufSize int) *PoolSync {
	return &PoolSync{
		pool: sync.Pool{
			New: func() interface{} {
				return NewBuffer(bufSize)
			},
		},
	}
}
