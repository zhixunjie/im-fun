package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// ReadWs read a proto from websocket connection.
func (proto *Proto) ReadWs(conn *websocket.Conn) (err error) {
	buf, err := conn.ReadMessage()
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
		proto.Body = buf[pack.HeaderLen:pack.BodyLen]
	} else {
		proto.Body = nil
	}
	return
}

// WriteWs write a proto to websocket connection.
func (proto *Proto) WriteWs(conn *websocket.Conn) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [header] + [body] ]

	// websocket header
	payloadLen := _rawHeaderSize + len(proto.Body)
	err = conn.WriteHeader(websocket.BinaryMessage, payloadLen)
	if err != nil {
		return err
	}

	// websocket payload
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	err = conn.WritePayload(codeHeader(proto, buf))
	if err != nil {
		return err
	}
	if proto.Body != nil {
		err = conn.WritePayload(proto.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteWsHeart write websocket heartbeat with room online.
func (proto *Proto) WriteWsHeart(conn *websocket.Conn, online int32) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [header] + [body] ]

	// websocket header
	payloadLen := _rawHeaderSize + _heartSize
	err = conn.WriteHeader(websocket.BinaryMessage, payloadLen)
	if err != nil {
		return err
	}

	// websocket payload
	buf := make([]byte, payloadLen) // TODO try to reduce GC
	buf = codeHeader(proto, buf)
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	err = conn.WritePayload(codeHeader(proto, buf))
	if err != nil {
		return err
	}

	return
}
