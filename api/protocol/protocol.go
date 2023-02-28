package protocol

import (
	"errors"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

// size
const (
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
)

// offset
const (
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
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

func parseHeader(proto *Proto, buf []byte) error {
	if len(buf) < _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	packLen := binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen := binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	proto.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	proto.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	proto.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
}
