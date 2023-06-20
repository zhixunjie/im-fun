package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

// 针对TCP连接的消息发送和接收

// ReadTCP read a proto from TCP reader.
func (proto *Proto) ReadTCP(reader *bufio.Reader) (err error) {
	var buf []byte

	// read n bytes
	if buf, err = reader.Pop(_rawHeaderSize); err != nil {
		return
	}

	// proto header
	pack, err := decodeHeaderFromBufToProto(buf, proto)
	if err != nil {
		return err
	}
	// proto body
	if pack.BodyLen > 0 {
		proto.Body, err = reader.Pop(int(pack.BodyLen))
	} else {
		proto.Body = nil
	}
	return nil
}

// WriteTCP write a proto to TCP writer.
func (proto *Proto) WriteTCP(writer *bufio.Writer) (err error) {
	// raw message（no header，send between service，only service job will send this kind of msg by now）
	if proto.Op == int32(OpBatchMsg) {
		_, err = writer.Write(proto.Body)
		return
	} else {
		// Peek：只需要把header的内存区peek出来即可
		var buf []byte
		if buf, err = writer.Peek(_rawHeaderSize); err != nil {
			return
		}

		// 1. proto header
		encodeHeaderFromProtoToBuf(proto, buf)
		// 2. proto body
		if proto.Body != nil {
			_, err = writer.Write(proto.Body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteTCPHeart write TCP heartbeat with room online.
func (proto *Proto) WriteTCPHeart(writer *bufio.Writer, online int32) (err error) {
	// Peek：一次性把整个数据包的内存区都Peek出来
	var buf []byte
	packLen := _rawHeaderSize + _heartSize
	if buf, err = writer.Peek(packLen); err != nil {
		return
	}

	// 1. proto header
	encodeHeaderFromProtoToBuf(proto, buf)
	// 2. proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)

	return nil
}
