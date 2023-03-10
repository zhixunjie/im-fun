package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// 针对WebSocket连接的消息发送和接收

// ReadWs read a proto from websocket connection.
func (proto *Proto) ReadWs(conn *websocket.Conn) (err error) {
	buf, err := conn.ReadMessage()
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
		proto.Body = buf[pack.HeaderLen:pack.PackLen]
	} else {
		proto.Body = nil
	}
	return
}

// WriteWs write a proto to websocket connection.
func (proto *Proto) WriteWs(conn *websocket.Conn) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [proto header] + [proto body] ]

	// websocket header
	payloadLen := _rawHeaderSize + len(proto.Body)
	err = conn.WriteHeader(websocket.BinaryMessage, payloadLen)
	if err != nil {
		return err
	}

	// websocket payload
	buf := make([]byte, _rawHeaderSize) // TODO try to reduce GC
	// proto header
	buf = codeProtoHeader(proto, buf)
	err = conn.WritePayload(buf)
	if err != nil {
		return err
	}
	// proto body
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
	// websocket payload = [ [proto header] + [proto body] ]

	// websocket header
	payloadLen := _rawHeaderSize + _heartSize
	err = conn.WriteHeader(websocket.BinaryMessage, payloadLen)
	if err != nil {
		return err
	}

	// websocket payload
	buf := make([]byte, payloadLen) // TODO try to reduce GC
	// proto header
	buf = codeProtoHeader(proto, buf)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	err = conn.WritePayload(codeProtoHeader(proto, buf))
	if err != nil {
		return err
	}

	return
}
