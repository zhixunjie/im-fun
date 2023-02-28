package protocol

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// ReadWebsocket read a proto from websocket connection.
func (proto *Proto) ReadWebsocket(conn *websocket.Conn) (err error) {
	buf, err := conn.ReadMessage()
	if err != nil {
		logrus.Errorf("err=%v,", err)
		return
	}
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

	// read body
	if bodyLen := int(packLen - int32(headerLen)); bodyLen > 0 {
		proto.Body = buf[headerLen:packLen]
	} else {
		proto.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (proto *Proto) WriteWebsocket(conn *websocket.Conn) (err error) {
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC

	// write header
	packLen := _rawHeaderSize + len(proto.Body)
	if err = conn.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)

	// write body
	if proto.Body != nil {
		err = conn.WriteBody(proto.Body)
	}
	return
}

// WriteWebsocketHeart write websocket heartbeat with room online.
func (proto *Proto) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
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
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(proto.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], proto.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], proto.Seq)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	// TODO .....
	return
}
