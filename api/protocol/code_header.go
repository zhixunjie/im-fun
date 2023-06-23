package protocol

import (
	"errors"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

const ProtoVersion = 1

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

// size
const (
	// package's size
	_packSize   = 4
	_headerSize = 2
	_verSize    = 2
	_opSize     = 4
	_seqSize    = 4
	// header's size
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize

	// heart's size
	_heartSize = 4

	// other size
	_maxPackSize = MaxBodySize + int32(_rawHeaderSize)
)

// offset
const (
	// package's size
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize

	// heart's offset
	_heartOffset = _seqOffset + _seqSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("pack total length not allow")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("pack header length not allow")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Op: int32(OpProtoReady)}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Op: int32(OpProtoFinish)}
)

type ProtoPack struct {
	PackLen   int32 // 整个数据包的长度
	HeaderLen int16 // 头部的长度
	BodyLen   int32 // 请求体的长度
}

// header编码：把proto的头信息，编码写入到buf的头
func encodeHeaderFromProtoToBuf(proto *Proto, buf []byte) {
	packLen := _rawHeaderSize + int32(len(proto.Body))
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
}

// header解码：把buf的头信息，解码到proto的头
func decodeHeaderFromBufToProto(buf []byte, proto *Proto) (ProtoPack, error) {
	var header ProtoPack
	if len(buf) < _rawHeaderSize {
		return header, ErrProtoHeaderLen
	}
	// parse header
	packLen := binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen := binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	proto.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	proto.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	proto.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	header.PackLen = packLen
	header.HeaderLen = headerLen
	header.BodyLen = packLen - int32(headerLen)
	// check message
	if packLen > _maxPackSize {
		return header, ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return header, ErrProtoHeaderLen
	}

	return header, nil
}
