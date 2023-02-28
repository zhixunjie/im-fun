package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// ReadWebsocket read a proto from websocket connection.
func (proto *Proto) ReadWebsocket(conn *websocket.Conn) (err error) {
	// read the whole message
	buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// parse header
	pack, err := unCode(proto, buf)
	if err != nil {
		return err
	}

	// read body
	if pack.BodyLen > 0 {
		proto.Body = buf[pack.HeaderLen:pack.BodyLen]
	} else {
		proto.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (proto *Proto) WriteWebsocket(conn *websocket.Conn) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [header] + [body] ]

	// websocket header
	packLen := _rawHeaderSize + len(proto.Body)
	if err = conn.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return err
	}

	// websocket payload
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = conn.WriteBody(code(proto, buf))
	if err != nil {
		return err
	}
	// write body
	if proto.Body != nil {
		err = conn.WriteBody(proto.Body)
		if err != nil {
			return err
		}
	}
	return nil
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
