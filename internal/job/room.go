package job

import (
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

var (
	ErrRoomFull = errors.New("room proto chan full")
)

// RoomJob 房间任务
type RoomJob struct {
	roomId  string               // 房间ID
	conf    *conf.Room           // 相关配置
	job     *Job                 // 依赖于Job对象
	protoCh chan *protocol.Proto // 有缓冲的Channel
}

func NewRoom(b *Job, roomId string) (r *RoomJob) {
	c := b.conf.Room
	r = &RoomJob{
		roomId:  roomId,
		conf:    c,
		job:     b,
		protoCh: make(chan *protocol.Proto, c.Batch*2),
	}
	go r.receiveFromCh(c.Batch, time.Duration(c.Interval))
	return
}

// SendToCh 向房间发送消息
func (r *RoomJob) SendToCh(msg []byte) error {
	logHead := fmt.Sprintf("SendToCh|msg=%s,", msg)

	var p = &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Seq:  gmodel.NewSeqId32(),
		Body: msg,
	}

	// try to put into channel, otherwise return ErrRoomFull.
	select {
	case r.protoCh <- p:
		logging.Infof(logHead + "send success")
		return nil
	default:
		logging.Infof(logHead + "send error")
		return ErrRoomFull
	}
}

func (r *RoomJob) receiveFromCh(accumulate int, interval time.Duration) {
	logHead := fmt.Sprintf("receiveFromCh|roowId=%v,", r.roomId)

	timerDuration := interval * 100
	timer := time.NewTimer(timerDuration)
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
		timer.Reset(timerDuration)
	}

	// 事件循环：有三种不同的策略
	logging.Infof(logHead + "start room's loop")
	for {
		select {
		case proto = <-r.protoCh:
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
		// 策略3：如果一个房间很久没有收到消息，那么就删除房间（释放内存）
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
