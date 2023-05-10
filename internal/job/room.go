package job

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"time"
)

var (
	ErrRoomFull    = errors.New("room proto chan full")
	roomReadyProto = new(protocol.Proto)
)

type Room struct {
	conf  *conf.Room
	job   *Job
	id    string               // 房间Id
	proto chan *protocol.Proto // 有缓冲的Channel
}

func NewRoom(job *Job, roomId string) (r *Room) {
	conf := job.conf.Room
	r = &Room{
		conf:  conf,
		id:    roomId,
		job:   job,
		proto: make(chan *protocol.Proto, conf.Batch*2),
	}
	go r.PopFromChannel(conf.Batch, time.Duration(conf.TimerDuration))
	return
}

func (r *Room) PushToChannel(msg []byte) (err error) {
	var p = &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Body: msg,
	}
	select {
	case r.proto <- p:
	default:
		err = ErrRoomFull
	}
	return
}

func (r *Room) PopFromChannel(batch int, timerDuration time.Duration) {
	logHead := "PopFromChannel"
	var (
		n    int
		last time.Time
		p    *protocol.Proto
		buf  = buffer.NewWriterSize(int(protocol.MaxBodySize))
	)
	logrus.Infof(logHead+"start roomId=%s", r.id)
	timer := time.AfterFunc(timerDuration, func() {
		select {
		case r.proto <- roomReadyProto:
		default:
		}
	})
	defer timer.Stop()

	// begin to traverse
	for {
		p = <-r.proto
		switch {
		case p == nil: //
			goto end
		case p != roomReadyProto:
			// merge buffer ignore error, always nil
			p.WriteTo(buf)
			if n++; n == 1 {
				last = time.Now()
				timer.Reset(timerDuration)
				continue
			} else if n < batch {
				if timerDuration > time.Since(last) {
					continue
				}
			}
		case n == 0:
			goto end
		}

		_ = r.job.PushUserRoom(0, r.id, buf.Buffer())
		// after push to room channel, renew a buffer, let old buffer gc
		// TODO use reset buffer will be better
		buf = buffer.NewWriterSize(buf.Size())
		n = 0
		if r.conf.Idle != 0 {
			timer.Reset(time.Duration(r.conf.Idle))
		} else {
			timer.Reset(time.Minute)
		}
	}
end:
	r.job.delRoom(r.id)
	logrus.Infof("room:%s goroutine exit", r.id)
}

func (j *Job) delRoom(roomID string) {
	j.roomsMutex.Lock()
	delete(j.rooms, roomID)
	j.roomsMutex.Unlock()
}

func (j *Job) getRoom(roomID string) *Room {
	j.roomsMutex.RLock()
	room, ok := j.rooms[roomID]
	j.roomsMutex.RUnlock()
	if !ok {
		j.roomsMutex.Lock()
		if room, ok = j.rooms[roomID]; !ok {
			room = NewRoom(roomID, j)
			j.rooms[roomID] = room
		}
		j.roomsMutex.Unlock()
		logrus.Infof("new a room:%s active:%d", roomID, len(j.rooms))
	}
	return room
}
