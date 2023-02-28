package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

// WriteTo write a proto to bytes writer.
func (proto *Proto) WriteTo(writer *buffer.Writer) {
	// writer header
	buf := writer.Peek(_rawHeaderSize)
	codeHeader(proto, buf)
	// writer body
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}

// ReadTCP read a proto from TCP reader.
func (proto *Proto) ReadTCP(reader *bufio.Reader) (err error) {
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = reader.ReadBytesN(buf)
	if err != nil {
		return err
	}
	// read header
	pack, err := unCodeHeader(proto, buf)
	if err != nil {
		return err
	}
	// read body
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
	// write raw message
	if proto.Op == int32(OpRaw) {
		_, err = writer.Write(proto.Body)
		return
	}
	// proto header
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	_, err = writer.Write(codeHeader(proto, buf))
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
func (proto *Proto) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
	// proto header
	packLen := _rawHeaderSize + _heartSize
	buf := make([]byte, packLen) // TODO try to reduce GC
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	_, err = wr.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
