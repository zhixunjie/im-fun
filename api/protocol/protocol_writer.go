package protocol

import "github.com/zhixunjie/im-fun/pkg/buffer"

// WriteTo write a proto to writer.
func (proto *Proto) WriteTo(writer *buffer.Writer) {
	// proto header
	buf := writer.Peek(_rawHeaderSize)
	// code proto to buf
	buf = codeProtoHeader(proto, buf)
	// proto body
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}
