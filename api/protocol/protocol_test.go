package protocol

import (
	"fmt"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	var proto = &Proto{
		Ver:  ProtoVersion,
		Op:   int32(OpAuth),
		Seq:  1,
		Body: []byte("i am jason"),
	}
	buf := make([]byte, _rawHeaderSize)
	codeProtoHeader(proto, buf)

	var proto1 Proto
	fmt.Println(unCodeProtoHeader(&proto1, buf))
}
