package channel

import (
	"errors"
	"github.com/zhixunjie/im-fun/api/protocol"
)

var (
	ErrRingEmpty = errors.New("ring buffer empty")
	ErrRingFull  = errors.New("ring buffer full")
)

// Ring use to get proto，reduce GC
// Every User Has A Channel and Every Channel Has a Ring
type Ring struct {
	// read
	rp   uint64
	max  uint64
	mask uint64

	// write
	wp   uint64
	data []protocol.Proto
}

func (r *Ring) Init(max uint64) {
	// max must be a number of 2^n
	if max&(max-1) != 0 {
		for max&(max-1) != 0 {
			max &= max - 1
		}
		max <<= 1
	}
	// make proto error
	r.data = make([]protocol.Proto, max)
	r.max = max
	r.mask = r.max - 1
}

// GetProtoCanRead 获取一个Proto（用于读取的Proto）
func (r *Ring) GetProtoCanRead() (proto *protocol.Proto, err error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	proto = &r.data[r.rp&r.mask] // r.rp % r.mask
	return
}

// GetProtoCanWrite 获取一个Proto（用于写入的Proto）
func (r *Ring) GetProtoCanWrite() (proto *protocol.Proto, err error) {
	// 超出一定范围就不能再写入了，前端的消息就发不过来了。
	if r.wp-r.rp >= r.max {
		return nil, ErrRingFull
	}
	proto = &r.data[r.wp&r.mask] // r.wp % r.mask
	return
}

// AdvReadPointer 向前推进读指针
func (r *Ring) AdvReadPointer() {
	r.rp++
}

// AdvWritePointer 向前推进写指针
func (r *Ring) AdvWritePointer() {
	r.wp++
}

// ResetPointer 重置指针
func (r *Ring) ResetPointer() {
	r.rp = 0
	r.wp = 0
}
