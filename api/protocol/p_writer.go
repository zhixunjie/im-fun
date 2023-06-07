package protocol

import "github.com/zhixunjie/im-fun/pkg/buffer"

// WriteTo
// write a proto to writer
func (proto *Proto) WriteTo(writer *buffer.BufferWriter) {
	buf := writer.Peek(_rawHeaderSize)

	// encode proto's header to buffer
	buf = encodeHeaderFromProtoToBuf(proto, buf)

	// proto body
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}
