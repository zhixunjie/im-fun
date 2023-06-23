package operation

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhixunjie/im-fun/benchmarks/client/tcp/model"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net"
	"sync/atomic"
	"time"
)

func Auth(rd *bufio.Reader, wr *bufio.Writer, userId int64) (err error) {
	logHead := fmt.Sprintf("auth|userId=%v,", userId)

	var authParams = &model.AuthParams{
		UserId:   userId,
		UserKey:  "random",
		RoomId:   "live://9999",
		Platform: "linux",
		Token:    "abcabcabcabc",
	}
	body, _ := json.Marshal(authParams)
	proto := &model.Proto{
		Ver:  1,
		Op:   model.OpAuth,
		Body: body,
	}

	// auth
	if err = WriteProto(wr, proto); err != nil {
		logging.Errorf(logHead+"write() error(%v)", err)
		return
	}
	logging.Infof(logHead+"auth req,authParams=%+v", authParams)

	// auth reply
	if proto, err = ReadProto(rd); err != nil {
		logging.Errorf(logHead+"read() error(%v)", err)
		return
	}
	logging.Infof(logHead+"auth reply,proto=%+v", proto)

	return
}

func Writer(ctx context.Context, seq *int32, wr *bufio.Writer, userId int64, quit chan bool) (err error) {
	logHead := fmt.Sprintf("Writer|userId=%v,", userId)
	proto := new(model.Proto)

	// deal with write
	for {
		select {
		case <-quit:
			return
		default:
		}

		// heartbeat
		proto.Op = model.OpHeartbeat
		proto.Seq = *seq
		proto.Body = nil
		*seq++

		// write proto
		if err = WriteProto(wr, proto); err != nil {
			logging.Errorf(logHead+"WriteProto() error(%v)", err)
			return
		}
		logging.Infof(logHead + "Write heartbeat success")
		time.Sleep(model.Heart)
	}
}

func Reader(ctx context.Context, conn net.Conn, rd *bufio.Reader, userId int64, quit chan bool) (err error) {
	logHead := fmt.Sprintf("Reader|userId=%v,", userId)
	proto := new(model.Proto)

	// deal with read
	for {
		// read proto
		if proto, err = ReadProto(rd); err != nil {
			logging.Errorf(logHead+"ReadProto() error(%v)", err)
			quit <- true
			return
		}

		// check operation
		switch proto.Op {
		case model.OpAuthReply:
			logging.Infof(logHead + "receive auth reply")
		case model.OpHeartbeatReply:
			logging.Infof(logHead+"receive heartbeat reply", userId)
			if err = conn.SetReadDeadline(time.Now().Add(model.Heart + 60*time.Second)); err != nil {
				logging.Errorf(logHead+"conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
		case model.OpBatchMsg:
			bodyLen := proto.BodyLen
			if bodyLen > 0 {
				var batchProto = new(model.Proto)
				for {
					batchProto, err = ReadProto(rd)
					if err != nil {
						logging.Errorf(logHead+"ReadProto() error(%v)", err)
						break
					}
					bodyLen -= batchProto.BodyLen
					logging.Infof(logHead+"receive msg,proto=%+v", batchProto)
					if bodyLen == 0 {
						break
					}
				}
				atomic.AddInt64(&mCount, 1)
			}
		default:
			logging.Infof(logHead+"receive unknown msg,proto=%+v", proto)
		}
	}
}
