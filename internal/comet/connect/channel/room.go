package channel

import (
	"sync"

	"github.com/zhixunjie/im-fun/api/protocol"
)

// Room is a room and store channel room info.
type Room struct {
	Id        string
	rLock     sync.RWMutex
	next      *Channel // linklist: ch1 -> ch2 -> ch3 （every connection has a channel）
	Online    int32    // dirty read is ok
	AllOnline int32
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
	r.rLock.Lock()
	if r.next != nil {
		r.next.Prev = ch
	}
	ch.Next = r.next
	ch.Prev = nil
	r.next = ch
	r.Online++
	r.rLock.Unlock()
	return
}

// DelChannel 从房间删除对象的Channel
func (r *Room) DelChannel(ch *Channel) {
	r.rLock.Lock()
	if ch.Next != nil { // if not tail in the linklist
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil { // if not head in the linklist
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	ch.Next = nil
	ch.Prev = nil
	r.Online--
	r.rLock.Unlock()
}

// PushToAllChan 把proto推送到房间中的所有Channel
func (r *Room) PushToAllChan(proto *protocol.Proto) {
	r.rLock.RLock()
	// if chan full，discard it
	for ch := r.next; ch != nil; ch = ch.Next {
		_ = ch.Push(proto)
	}
	r.rLock.RUnlock()
}

func (r *Room) Close() {
	r.rLock.RLock()
	// close channel one by one in the room
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
	r.rLock.RUnlock()
}

// OnlineNum 房间的在线人数
func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
