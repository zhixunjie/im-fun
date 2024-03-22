package protocol

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/encoding/binary"
	"github.com/zhixunjie/im-fun/pkg/websocket"
)

// TCP连接的读写操作
// 能够同时适配TCP与WebSocket

type ConnReadWriter interface {
	ReadProto(proto *Proto) (err error)
	WriteProto(proto *Proto) (err error)
	WriteProtoHeart(proto *Proto, online int32) (err error)
	Flush() error
}

type TcpConnReaderWriter struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTcpConnReaderWriter(reader *bufio.Reader, writer *bufio.Writer) *TcpConnReaderWriter {
	r := new(TcpConnReaderWriter)
	r.reader = reader
	r.writer = writer
	return r
}

// ReadProto read a proto from TCP connection(reader/fd)
func (r *TcpConnReaderWriter) ReadProto(proto *Proto) (err error) {
	var buf []byte

	// 1. read n bytes
	reader := r.reader
	if buf, err = reader.Pop(_rawHeaderSize); err != nil {
		return
	}

	// 2. 把buf的头信息，解码到proto的头
	pack, err := decodeHeaderFromBufToProto(buf, proto)
	if err != nil {
		return
	}
	// 3. 把buf的Body信息，放入到proto
	if pack.BodyLen > 0 {
		proto.Body, err = reader.Pop(int(pack.BodyLen)) // 从reader中Pop出N个字节
	} else {
		proto.Body = nil
	}
	return
}

// WriteProto writer a proto to TCP connection(writer/fd)
func (r *TcpConnReaderWriter) WriteProto(proto *Proto) (err error) {
	writer := r.writer

	// 1. Peek：只需要把header的内存区peek出来即可
	var buf []byte
	if buf, err = writer.Peek(_rawHeaderSize); err != nil {
		return
	}
	// 2. 把proto的头信息，编码写入到buf的头
	packLen := _rawHeaderSize + int32(len(proto.Body))
	encodeHeaderFromProtoToBuf(packLen, proto, buf)
	// 3. 把proto的Body信息，编码写入到buf中
	if proto.Body != nil {
		_, err = writer.Write(proto.Body)
		if err != nil {
			return
		}
	}

	return
}

// WriteProtoHeart write TCP heartbeat with room online
func (r *TcpConnReaderWriter) WriteProtoHeart(proto *Proto, online int32) (err error) {
	writer := r.writer

	// 1. Peek：一次性把整个数据包的内存区都Peek出来
	var buf []byte
	packLen := _rawHeaderSize + _heartSize
	if buf, err = writer.Peek(packLen); err != nil {
		return
	}
	// 2. 把proto的头信息，编码写入到buf的头
	encodeHeaderFromProtoToBuf(int32(packLen), proto, buf)
	// 3. 把proto的Body信息，编码写入到buf中
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)

	return
}

func (r *TcpConnReaderWriter) Flush() error {
	return r.writer.Flush()
}

type WsConnReaderWriter struct {
	conn *websocket.Conn
}

func NewWsConnReaderWriter(conn *websocket.Conn) *WsConnReaderWriter {
	r := new(WsConnReaderWriter)
	r.conn = conn
	return r
}

// ReadProto read a proto from WebSocket connection(reader/fd)
func (r *WsConnReaderWriter) ReadProto(proto *Proto) (err error) {
	conn := r.conn

	// 1. read message from websocket
	buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// 2. 把buf的头信息，解码到proto的头
	pack, err := decodeHeaderFromBufToProto(buf, proto)
	if err != nil {
		return
	}
	// 3. 把buf的Body信息，放入到proto
	if pack.BodyLen > 0 {
		proto.Body = buf[pack.HeaderLen:pack.PackLen] // 从buf中切分出N个字节
	} else {
		proto.Body = nil
	}
	return
}

// WriteProto writer a proto to WebSocket connection(writer/fd)
func (r *WsConnReaderWriter) WriteProto(proto *Proto) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [proto header] + [proto body] ]

	// 1. write websocket header
	conn := r.conn
	packLen := _rawHeaderSize + len(proto.Body)
	err = conn.WriteHeader(websocket.BinaryMessage, packLen)
	if err != nil {
		return err
	}

	// 2. write websocket payload
	{
		// 2.1 Peek：只需要把header的内存区peek出来即可
		var buf []byte
		if buf, err = conn.Peek(_rawHeaderSize); err != nil {
			return
		}
		// 2.2 把proto的头信息，编码写入到buf的头
		encodeHeaderFromProtoToBuf(int32(packLen), proto, buf)
		// 2.3 把proto的Body信息，编码写入到buf中
		if proto.Body != nil {
			err = conn.WritePayload(proto.Body)
			if err != nil {
				return
			}
		}
	}

	return
}

// WriteProtoHeart write websocket heartbeat with room online
func (r *WsConnReaderWriter) WriteProtoHeart(proto *Proto, online int32) (err error) {
	// format:
	// pack = [ [websocket header] + [websocket payload]]
	// websocket payload = [ [proto header] + [proto body] ]

	// 1. part1: websocket header
	conn := r.conn
	packLen := _rawHeaderSize + _heartSize // websocket payload
	err = conn.WriteHeader(websocket.BinaryMessage, packLen)
	if err != nil {
		return err
	}

	// 2. part2: websocket payload
	{
		// 2.1 Peek：一次性把整个数据包的内存区都Peek出来
		var buf []byte
		if buf, err = conn.Peek(packLen); err != nil {
			return
		}
		// 2.2 把proto的头信息，编码写入到buf的头
		encodeHeaderFromProtoToBuf(int32(packLen), proto, buf)
		// 2.3 把proto的Body信息，编码写入到buf中
		binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	}

	return
}

func (r *WsConnReaderWriter) Flush() error {
	return r.conn.Flush()
}

// WriteProtoToWriter write a proto to writer(带有buf的writer)
func WriteProtoToWriter(proto *Proto, writer *bytes.BufferWriter) {
	// 1. Peek：只需要把header的内存区peek出来即可
	buf := writer.Peek(_rawHeaderSize)
	// 2. 把proto的头信息，编码写入到buf的头
	packLen := _rawHeaderSize + int32(len(proto.Body))
	encodeHeaderFromProtoToBuf(packLen, proto, buf)
	// 3. 把proto的Body信息，编码写入到buf中
	if proto.Body != nil {
		writer.Write(proto.Body)
	}
}
