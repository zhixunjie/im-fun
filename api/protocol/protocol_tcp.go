package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

// 针对TCP连接的消息发送和接收

// ReadTCP read a proto from TCP reader.
func (proto *Proto) ReadTCP(reader *bufio.Reader) (err error) {
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = reader.ReadBytesN(buf)
	if err != nil {
		return err
	}
	// proto header
	pack, err := unCodeProtoHeader(proto, buf)
	if err != nil {
		return err
	}
	// proto body
	if pack.BodyLen > 0 {
		proto.Body = make([]byte, pack.BodyLen) // TODO try to reduce GC
		err = reader.ReadBytesN(proto.Body)
		if err != nil {
			return err
		}
	} else {
		proto.Body = nil
	}
	return nil
}

// WriteTCP write a proto to TCP writer.
func (proto *Proto) WriteTCP(writer *bufio.Writer) (err error) {
	// raw message（no header，send between service，only service job will send this kind of msg by now）
	if proto.Op == int32(OpRaw) {
		_, err = writer.Write(proto.Body)
		return
	}
	// proto header
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	buf = codeProtoHeader(proto, buf)
	_, err = writer.Write(buf)
	if err != nil {
		return err
	}
	// proto body
	if proto.Body != nil {
		_, err = writer.Write(proto.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteTCPHeart write TCP heartbeat with room online.
func (proto *Proto) WriteTCPHeart(writer *bufio.Writer, online int32) (err error) {
	// proto header
	packLen := _rawHeaderSize + _heartSize
	buf := make([]byte, packLen) // TODO try to reduce GC
	buf = codeProtoHeader(proto, buf)

	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	_, err = writer.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
