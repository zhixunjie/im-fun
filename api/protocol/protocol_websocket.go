package protocol

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// ReadWebsocket read a proto from websocket connection.
func (p *Proto) ReadWebsocket(conn *websocket.Conn) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
	)
	buf, err := conn.ReadMessage()
	if err != nil {
		logrus.Errorf("err=%v,", err)
		return
	}
	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}

	// read body
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Proto) WriteWebsocket(conn *websocket.Conn) (err error) {
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC

	// write header
	packLen := _rawHeaderSize + len(p.Body)
	if err = conn.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)

	// write body
	if p.Body != nil {
		err = conn.WriteBody(p.Body)
	}
	return
}

// WriteWebsocketHeart write websocket heartbeat with room online.
func (p *Proto) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
	packLen := _rawHeaderSize + _heartSize
	buf := make([]byte, packLen) // TODO try to reduce GC

	// write header
	err = wr.WriteHeader(websocket.BinaryMessage, packLen)
	if err != nil {
		return
	}
	// proto header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	// TODO .....
	return
}
