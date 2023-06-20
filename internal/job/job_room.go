package job

import (
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

var (
	ErrRoomFull = errors.New("room proto chan full")
)

type Room struct {
	conf  *conf.Room
	job   *Job
	id    string               // 房间Id
	proto chan *protocol.Proto // 有缓冲的Channel
}

func NewRoom(job *Job, roomId string) (r *Room) {
	c := job.conf.Room
	r = &Room{
		conf:  c,
		id:    roomId,
		job:   job,
		proto: make(chan *protocol.Proto, c.Batch*2),
	}
	go r.Receive(c.Batch, time.Duration(c.Duration))
	return
}

// Send 向房间发送消息
func (r *Room) Send(msg []byte) error {
	var p = &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Body: msg,
	}

	// try to put into channel, otherwise return ErrRoomFull
	select {
	case r.proto <- p:
		return nil
	default:
		return ErrRoomFull
	}
}

func (r *Room) Receive(batch int, duration time.Duration) {
	logHead := fmt.Sprintf("Receive|roowId=%v,", r.id)

	ticker := time.NewTicker(duration)
	timerDuration := duration * 5000
	timer := time.NewTicker(timerDuration)
	defer timer.Stop()

	var writer = bytes.NewWriterSize(int(protocol.MaxBodySize))
	var n int
	var proto *protocol.Proto

	logging.Infof(logHead + "new room")
	for {
		select {
		// 策略1：每个一段时间发送一次群消息
		case <-ticker.C:
			logging.Infof(logHead + "ticker.C")
			break
		// 策略2：累积到一定数目后发送一次群消息
		case proto = <-r.proto:
			if proto != nil {
				logging.Infof(logHead+"get proto=%v,n=%v", proto, n)
				protocol.WriteProtoToWriter(proto, writer)
				if n++; n >= batch {
					break
				}
			}
		// 策略3：如果很久没有收到消息，那么就删除房间（释放内存）
		case <-timer.C:
			goto end
		}
		content := writer.Buffer()
		if len(content) == 0 {
			continue
		}
		_ = r.job.SendToRoom(0, r.id, content)
		n = 0
		writer.Reset()
		timer.Reset(timerDuration)
	}
end:
	r.job.DelRoom(r.id)
	logging.Infof(logHead + "delete room")
}

func (job *Job) DelRoom(roomId string) {
	job.rwMutex.Lock()
	delete(job.rooms, roomId)
	job.rwMutex.Unlock()
}

func (job *Job) CreateOrGetRoom(roomId string) *Room {
	job.rwMutex.RLock()
	room, ok := job.rooms[roomId]
	job.rwMutex.RUnlock()
	if !ok {
		job.rwMutex.Lock()
		if room, ok = job.rooms[roomId]; !ok {
			room = NewRoom(job, roomId)
			job.rooms[roomId] = room
		}
		job.rwMutex.Unlock()
		logging.Infof("create a room=%s,active=%d", roomId, len(job.rooms))
	} else {
		logging.Infof("get a room=%s,active=%d", roomId, len(job.rooms))
	}
	return room
}
