package job

import (
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

var (
	ErrRoomFull = errors.New("room proto chan full")
)

type RoomJob struct {
	conf   *conf.Room
	job    *Job
	roomId string               // 房间Id
	proto  chan *protocol.Proto // 有缓冲的Channel
}

func NewRoom(job *Job, roomId string) (r *RoomJob) {
	c := job.conf.Room
	r = &RoomJob{
		conf:   c,
		roomId: roomId,
		job:    job,
		proto:  make(chan *protocol.Proto, c.Batch*2),
	}
	go r.receiveFromCh(c.Batch, time.Duration(c.Duration))
	return
}

// SendToCh 向房间发送消息
func (r *RoomJob) SendToCh(msg []byte) error {
	var p = &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Seq:  int32(gen_id.SeqId()),
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

func (r *RoomJob) receiveFromCh(accumulate int, interval time.Duration) {
	logHead := fmt.Sprintf("receive|roowId=%v,", r.roomId)

	duration := interval * 100
	timer := time.NewTicker(duration)
	defer timer.Stop()

	var writer = bytes.NewWriterSize(int(protocol.MaxBodySize))
	var counter int
	var proto *protocol.Proto
	last := time.Now()

	fn := func(from string) {
		content := writer.Buffer()
		if len(content) == 0 {
			return
		}

		// send room msg
		_ = r.job.SendToRoom(0, r.roomId, content)

		// reset
		counter = 0
		writer.Reset()
		last = time.Now()
		timer.Reset(duration)
	}

	logging.Infof(logHead + "create room")
	for {
		select {
		case proto = <-r.proto:
			if proto != nil {
				logging.Infof(logHead+"get proto=%v,n=%v", proto, counter)
				protocol.WriteProtoToWriter(proto, writer)
				counter++
				// 策略1：累积到一定数目后发送一次群消息
				// if counter equal the value, then send msg to room
				if counter >= accumulate {
					fn("accumulate")
				} else {
					// 策略2：每隔一段时间发送一次群消息
					// if did not send since last time, then send msg to room
					if time.Since(last) >= interval && counter > 0 {
						fn("interval")
					}
				}
			}
		// 策略3：如果很久没有收到消息，那么就删除房间（释放内存）
		case <-timer.C:
			goto end
		}
	}
end:
	logging.Infof(logHead + "delete room")
	r.job.DelRoom(r.roomId)
}

// Job's Operation about RoomJob

func (b *Job) CreateOrGetRoom(roomId string) *RoomJob {
	b.rwMutex.RLock()
	room, ok := b.roomJobs[roomId]
	b.rwMutex.RUnlock()
	if !ok {
		b.rwMutex.Lock()
		if room, ok = b.roomJobs[roomId]; !ok {
			room = NewRoom(b, roomId)
			b.roomJobs[roomId] = room
		}
		b.rwMutex.Unlock()
		logging.Infof("create a room=%s,active=%d", roomId, len(b.roomJobs))
	} else {
		logging.Infof("get a room=%s,active=%d", roomId, len(b.roomJobs))
	}
	return room
}

func (b *Job) DelRoom(roomId string) {
	b.rwMutex.Lock()
	delete(b.roomJobs, roomId)
	b.rwMutex.Unlock()
}
