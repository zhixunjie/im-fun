package protocol

import (
	"bytes"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"testing"
)

func TestTcpConnReaderWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	reader := bufio.NewReader(buf)
	writer := bufio.NewWriter(buf)

	io := NewTcpConnReaderWriter(reader, writer)
	var proto = &Proto{
		Ver:  ProtoVersion,
		Op:   int32(OpBatchMsg),
		Seq:  1,
		Body: []byte("i am jason"),
	}

	err := io.WriteProto(proto)
	if err != nil {
		logging.Infof("io.WriteProto err=%v", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		logging.Infof("io.Flush err=%v", err)
		return
	}
	newProto := new(Proto)
	err = io.ReadProto(newProto)
	if err != nil {
		logging.Infof("io.ReadProto err=%v", err)
		return
	}
	fmt.Printf("proto=%+v\n", newProto)
}
