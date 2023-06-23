package protocol

import (
	"fmt"
	"testing"
)

func TestHeaderEncodeAndDecode(t *testing.T) {
	var proto = &Proto{
		Ver:  ProtoVersion,
		Op:   int32(OpAuth),
		Seq:  1,
		Body: []byte("i am jason"),
	}
	buf := make([]byte, _rawHeaderSize)
	packLen := _rawHeaderSize + int32(len(proto.Body))
	encodeHeaderFromProtoToBuf(packLen, proto, buf)

	var proto1 Proto
	r1, r2 := decodeHeaderFromBufToProto(buf, &proto1)
	fmt.Printf("%+v,err=%v\n", r1, r2)
}
