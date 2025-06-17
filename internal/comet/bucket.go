package comet

import (
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"go.uber.org/atomic"
	"sync"
)

// Bucket global value, use bucket to manage all the channel（all TCP connection）
type Bucket struct {
	conf   *conf.Bucket
	rwLock sync.RWMutex                // protect the channels for chs
	chs    map[string]*channel.Channel // map：model.UserInfo.TcpSessionId => GetChannel

	// room
	routineCounter atomic.Uint64
	rooms          map[string]*channel.Room // map: RoomId => Room
	routines       []chan *pb.SendToRoomReq // deal with proto to room

	ipCount map[string]int32
}

func NewBucket(conf *conf.Bucket) *Bucket {
	b := &Bucket{
		conf:     conf,
		chs:      make(map[string]*channel.Channel, conf.InitSizeChannelMap),
		rooms:    make(map[string]*channel.Room, conf.InitSizeRoomMap),
		routines: make([]chan *pb.SendToRoomReq, conf.RoutineAmount),
		ipCount:  make(map[string]int32),
	}

	// init routines：处理房间的广播事件
	for i := 0; i < conf.RoutineAmount; i++ {
		b.routines[i] = make(chan *pb.SendToRoomReq, conf.RoutineChannelSize)
		go b.ProcessProtoToRoom(i)
	}
	return b
}

func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

func (b *Bucket) RoomCount() int {
	return len(b.rooms)
}

// RoomsCount get all room id where online number > 0.
func (b *Bucket) RoomsCount() (res map[string]int32) {
	var (
		roomID string
		room   *channel.Room
	)
	b.rwLock.RLock()
	res = make(map[string]int32)
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = room.Online
		}
	}
	b.rwLock.RUnlock()
	return
}

// ChangeRoom change ro room
func (b *Bucket) ChangeRoom(newRoomId string, ch *channel.Channel) (err error) {
	logging.Errorf("traceId=%v] ChangeRoom|", ch.TraceId)

	var newRoom *channel.Room
	var ok bool
	var oldRoom = ch.Room

	// reset room
	if newRoomId == "" {
		if oldRoom != nil {
			oldRoom.DelChannel(ch)
			b.DelRoomById(oldRoom)
		}
		ch.Room = nil
		return
	}

	// get room info
	b.rwLock.Lock()
	if newRoom, ok = b.rooms[newRoomId]; !ok {
		newRoom = channel.NewRoom(newRoomId)
		b.rooms[newRoomId] = newRoom
	}
	b.rwLock.Unlock()

	// change room
	if oldRoom != nil {
		oldRoom.DelChannel(ch)
		b.DelRoomById(oldRoom)
	}

	err = newRoom.PutChannel(ch)
	if err != nil {
		return
	}
	ch.Room = newRoom
	return
}

func (b *Bucket) Put(ch *channel.Channel) (err error) {
	logHead := "Put|"

	var room *channel.Room
	var ok bool
	userInfo := ch.UserInfo
	tcpSessionId := userInfo.Connect.SessionId
	roomId := userInfo.RoomId

	b.rwLock.Lock()
	// close old channel
	if oldCh := b.chs[tcpSessionId]; oldCh != nil {
		oldCh.SendFinish(logHead)
	}
	// set new channel
	b.chs[tcpSessionId] = ch
	if roomId != "" {
		if room, ok = b.rooms[roomId]; !ok {
			room = channel.NewRoom(roomId)
			b.rooms[roomId] = room
		}
		ch.Room = room
	}
	b.ipCount[userInfo.ClientIP]++
	b.rwLock.Unlock()

	// put channel to the room
	if room != nil {
		err = room.PutChannel(ch)
	}
	return
}

// DelChannel 删除一个用户的Channel
func (b *Bucket) DelChannel(currCh *channel.Channel) {
	logging.Infof("[traceId=%v] DelChannel|", currCh.TraceId)
	userInfo := currCh.UserInfo
	tcpSessionId := userInfo.Connect.SessionId

	b.rwLock.Lock()
	if ch, ok := b.chs[tcpSessionId]; ok {
		// delete channel
		if ch == currCh {
			delete(b.chs, tcpSessionId)
		}
		// update ip counter
		if b.ipCount[userInfo.ClientIP] > 1 {
			b.ipCount[userInfo.ClientIP]--
		} else {
			delete(b.ipCount, userInfo.ClientIP)
		}
	}
	b.rwLock.Unlock()

	// delete channel in the room
	room := currCh.Room
	if room != nil {
		room.DelChannel(currCh)
		// if empty room, must delete from bucket
		b.DelRoomById(room)
	}
}

func (b *Bucket) GetChannel(key string) (ch *channel.Channel) {
	b.rwLock.RLock()
	ch = b.chs[key]
	b.rwLock.RUnlock()
	return
}

// GetRoomById 通过房间ID，获取房间
func (b *Bucket) GetRoomById(roomId string) (room *channel.Room) {
	b.rwLock.RLock()
	room = b.rooms[roomId]
	b.rwLock.RUnlock()
	return
}

// DelRoomById 通过房间ID，删除房间
func (b *Bucket) DelRoomById(room *channel.Room) {
	b.rwLock.Lock()
	delete(b.rooms, room.Id)
	b.rwLock.Unlock()
	room.Close()
}
