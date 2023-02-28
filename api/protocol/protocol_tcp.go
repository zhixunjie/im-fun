package protocol

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

// WriteTo write a proto to bytes writer.
func (proto *Proto) WriteTo(writer *buffer.Writer) {
	// writer header
	packLen := _rawHeaderSize + int32(len(proto.Body))
	buf := writer.Peek(_rawHeaderSize)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
	// writer body
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}

// ReadTCP read a proto from TCP reader.
func (proto *Proto) ReadTCP(reader *bufio.Reader) (err error) {
	// read header
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = reader.ReadBytesN(buf)
	if err != nil {
		return err
	}
	// parse header
	packLen := binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen := binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	proto.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	proto.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	proto.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	// read body
	bodyLen := int(packLen - int32(headerLen))
	if bodyLen > 0 {
		proto.Body = make([]byte, bodyLen) // TODO try to reduce GC
		err = reader.ReadBytesN(proto.Body)
		if err != nil {
			logrus.Errorf("err=%v,", err)
			return err
		}
	} else {
		proto.Body = nil
	}
	return nil
}

// WriteTCP write a proto to TCP writer.
func (proto *Proto) WriteTCP(writer *bufio.Writer) (err error) {
	if proto.Op == int32(OpRaw) {
		_, err = writer.Write(proto.Body)
		return
	}
	// write header
	packLen := _rawHeaderSize + int32(len(proto.Body))
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
	_, err = writer.Write(buf)
	if err != nil {
		logrus.Errorf("err=%v,", err)
		return err
	}
	// write body
	if proto.Body != nil {
		_, err = writer.Write(proto.Body)
		if err != nil {
			logrus.Errorf("err=%v,", err)
			return err
		}
	}
	return nil
}

// WriteTCPHeart write TCP heartbeat with room online.
func (proto *Proto) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
	// write header
	packLen := _rawHeaderSize + _heartSize
	buf := make([]byte, packLen) // TODO try to reduce GC
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
	// write body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	_, err = wr.Write(buf)
	if err != nil {
		logrus.Errorf("err=%v,", err)
	}
	return
}
