package operation

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/benchmarks/client/tcp/model"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"io"
	"net"
	"time"
)

func Auth(rd *bufio.Reader, wr *bufio.Writer, userId uint64, token string) (err error) {
	logHead := fmt.Sprintf("auth|userId=%v,", userId)
	authParams := &pb.AuthParams{
		UniId:    cast.ToString(userId),
		Token:    token,
		RoomId:   "live://9999",
		Platform: pb.Platform_Platform_PC,
	}
	body, _ := json.Marshal(authParams)
	proto := &model.Proto{
		Ver:  1,
		Op:   model.OpAuth,
		Body: body,
	}

	// auth
	if err = WriteProto(wr, proto); err != nil {
		logging.Errorf(logHead+"write() error=%v", err)
		return
	}
	logging.Infof(logHead+"auth req,authParams=%+v", authParams)

	// auth reply
	if err = ReadProto(rd, proto); err != nil {
		logging.Errorf(logHead+"ReadProto() error=%v", err)
		return
	}
	PrintProto(logHead+"receive reply auth", proto)

	return
}

func Writer(ctx context.Context, seq *int32, wr *bufio.Writer, userId uint64, quit chan bool) (err error) {
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
			logging.Errorf(logHead+"WriteProto() error=%v", err)
			return
		}
		logging.Infof(logHead + "Write heartbeat success")
		time.Sleep(model.Heart)
	}
}

func Reader(ctx context.Context, conn net.Conn, rd *bufio.Reader, userId uint64, quit chan bool) (err error) {
	logHead := fmt.Sprintf("Reader|userId=%v,", userId)

	// deal with read
	for {
		// read proto
		proto := new(model.Proto)
		if err = ReadProto(rd, proto); err != nil {
			if err == io.EOF {
				continue
			}
			logging.Errorf(logHead+"ReadProto() error=%v", err)
			quit <- true
			return
		}

		// check operation
		switch proto.Op {
		case model.OpAuthReply:
			PrintProto(logHead+"receive reply auth", proto)
		case model.OpHeartbeatReply:
			PrintProto(logHead+"receive reply heartbeat", proto)
			// set read deadline
			if err = conn.SetReadDeadline(time.Now().Add(model.Heart + 60*time.Second)); err != nil {
				logging.Errorf(logHead+"conn.SetReadDeadline() error=%v", err)
				quit <- true
				return
			}
		case model.OpBatchMsg:
			bodyLen := proto.BodyLen
			if bodyLen > 0 {
				buf := bufio.NewReader(bytes.NewReader(proto.Body))
				var batchProto = new(model.Proto)

				// 开始遍历batch消息的body
				for i := 0; i < int(bodyLen); i += int(batchProto.PackLen) {
					err = ReadProto(buf, batchProto)
					if err != nil {
						logging.Errorf(logHead+"ReadProto() error=%v", err)
						break
					}
					PrintProto(logHead+"receive msg", batchProto)
				}
				mCount.Add(1)
			}
		default:
			PrintProto(logHead+"receive unknown msg", proto)
		}
	}
}

func PrintProto(logHead string, proto *model.Proto) {
	logging.Infof(logHead+"(PackLen=%v,HeaderLen=%v,Ver=%v,Seq=%v,Body=%s)",
		proto.PackLen, proto.HeaderLen, proto.Ver, proto.Seq, proto.Body)
}
