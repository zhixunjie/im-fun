package protocol

import (
	"fmt"
	"log"
	"testing"
)

func TestHeaderEncodeAndDecode(t *testing.T) {
	var proto = &Proto{
		Ver:  ProtoVersion,
		Op:   int32(OpAuthReq),
		Seq:  1,
		Body: []byte("i am jason"),
	}
	buf := make([]byte, _rawHeaderSize)
	packLen := _rawHeaderSize + int32(len(proto.Body))
	encodeHeaderFromProtoToBuf(packLen, proto, buf)

	var proto1 Proto
	result, err := decodeHeaderFromBufToProto(buf, &proto1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("result=%+v\n", result)
}
