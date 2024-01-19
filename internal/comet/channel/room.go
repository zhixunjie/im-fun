package channel

import (
	"sync"

	"github.com/zhixunjie/im-fun/api/protocol"
)

// Room is a room and store channel room info.
type Room struct {
	Id        string
	rwLock    sync.RWMutex
	linklist  *Channel // linklist: ch1 -> ch2 -> ch3 （every connection has a channel）
	Online    int32    // 单台机器，单个房间的在线人数：dirty read is ok
	AllOnline int32    // 单台机器，所有房间的在线人数
}

// NewRoom new a room struct, store channel room info.
func NewRoom(id string) (r *Room) {
	return &Room{
		Id: id,
	}
}

// PutChannel 把Channel放到房间中
// insert to the head of the linklist
func (r *Room) PutChannel(ch *Channel) (err error) {
	r.rwLock.Lock()
	if r.linklist != nil {
		r.linklist.Prev = ch
	}
	ch.Next = r.linklist
	ch.Prev = nil
	r.linklist = ch
	r.Online++
	r.rwLock.Unlock()
	return
}

// DelChannel 从房间删除对象的Channel
func (r *Room) DelChannel(ch *Channel) {
	r.rwLock.Lock()
	if ch.Next != nil { // if not tail in the linklist
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil { // if not head in the linklist
		ch.Prev.Next = ch.Next
	} else {
		r.linklist = ch.Next
	}
	ch.Next = nil
	ch.Prev = nil
	r.Online--
	r.rwLock.Unlock()
}

// SendToAllChan 把proto推送到房间中的所有Channel
func (r *Room) SendToAllChan(proto *protocol.Proto) {
	r.rwLock.RLock()
	// if chan full，discard it
	for ch := r.linklist; ch != nil; ch = ch.Next {
		_ = ch.Push(proto)
	}
	r.rwLock.RUnlock()
}

func (r *Room) Close() {
	r.rwLock.RLock()
	// close channel one by one in the room
	for ch := r.linklist; ch != nil; ch = ch.Next {
		ch.SendFinish()
	}
	r.rwLock.RUnlock()
}

// OnlineNum 房间的在线人数
func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
