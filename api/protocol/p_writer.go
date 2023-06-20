package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
)

// WriteTo
// write a proto to writer
func (proto *Proto) WriteTo(writer *bytes.BufferWriter) {
	// 1. Peek：只需要把header的内存区peek出来即可
	buf := writer.Peek(_rawHeaderSize)

	// 2. 把proto的头信息，编码写入到buf的头
	encodeHeaderFromProtoToBuf(proto, buf)

	// 3. 把proto的Body信息，编码写入到buf中
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}
