package protocol

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/buffer"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Op: int32(OpProtoReady)}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Op: int32(OpProtoFinish)}
)

// WriteTo write a proto to bytes writer.
func (p *Proto) WriteTo(writer *buffer.Writer) {
	packLen := _rawHeaderSize + int32(len(p.Body))

	// writer header
	buf := writer.Peek(_rawHeaderSize)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)

	// writer body
	if p.Body != nil {
		writer.Write(p.Body)
	}
}

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(reader *bufio.Reader) (err error) {
	// 1. read header
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = reader.ReadBytesN(buf)
	if err != nil {
		return
	}

	// 2. parse header
	packLen := binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen := binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])

	// 3. check data
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}

	// 4. read body
	bodyLen := int(packLen - int32(headerLen))
	if bodyLen > 0 {
		p.Body = make([]byte, bodyLen) // TODO try to reduce GC
		err = reader.ReadBytesN(p.Body)
		if err != nil {
			logrus.Errorf("err=%v,", err)
		}
	} else {
		p.Body = nil
	}
	return
}

// WriteTCP write a proto to TCP writer.
func (p *Proto) WriteTCP(writer *bufio.Writer) (err error) {
	if p.Op == int32(OpRaw) {
		_, err = writer.Write(p.Body)
		return
	} else {
		buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC

		// write header
		packLen := _rawHeaderSize + int32(len(p.Body))
		binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
		binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
		binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
		binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
		binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
		_, err = writer.Write(buf)

		// write body
		if p.Body != nil {
			_, err = writer.Write(p.Body)
			if err != nil {
				logrus.Errorf("err=%v,", err)
			}
		}
		return
	}

}

// WriteTCPHeart write TCP heartbeat with room online.
func (p *Proto) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
	packLen := _rawHeaderSize + _heartSize
	buf := make([]byte, packLen) // TODO try to reduce GC

	// write header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	// write body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	_, err = wr.Write(buf)
	if err != nil {
		logrus.Errorf("err=%v,", err)
	}
	return
}
