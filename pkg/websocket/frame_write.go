package websocket

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
)

// WriteMessage
// write a message（跟read反着来就好）
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.2
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.5
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.6
func (c *Conn) WriteMessage(msgType int, msg []byte) (err error) {
	// 1. write header
	err = c.WriteHeader(msgType, len(msg))
	if err != nil {
		return
	}
	// 2. write payload
	err = c.WritePayload(msg)
	return
}

// WriteHeader
// write header frame.
func (c *Conn) WriteHeader(opcode int, length int) error {
	// 提前预留2个字节（用于后续的写入）
	var err error
	var h []byte
	if h, err = c.writer.Peek(2); err != nil {
		return err
	}

	// 1. set first byte(8bit)
	firstByte := &h[0]
	*firstByte = byte(0)
	*firstByte = *firstByte | fin          // fin
	*firstByte = *firstByte | byte(opcode) // opcode

	// //////////////////////////////
	// 2. set second byte(8bit)
	secondByte := &h[1]
	*secondByte = byte(0)
	switch {
	case length < 126: // 7 bits
		// write second byte
		*secondByte = *secondByte | byte(length)
	case length < 65536: // 16 bits
		// write second byte
		*secondByte = *secondByte | 126
		// write other byte
		var otherByte []byte
		if otherByte, err = c.writer.Peek(2); err != nil {
			return err
		}
		binary.BigEndian.PutUint16(otherByte, uint16(length))
	default: // 64 bits
		// write second byte
		*secondByte = *secondByte | 127

		// write other byte
		var otherByte []byte
		if otherByte, err = c.writer.Peek(8); err != nil {
			return err
		}
		binary.BigEndian.PutUint64(otherByte, uint64(length))
	}
	return nil
}

// WritePayload
// write payload
func (c *Conn) WritePayload(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.writer.Write(b)
		if err != nil {
			logrus.Errorf("Write err=%v", err)
			return
		}
	}
	return
}

func (c *Conn) Peek(n int) ([]byte, error) {
	return c.writer.Peek(n)
}

// Flush
// flush writer buffer
func (c *Conn) Flush() error {
	return c.writer.Flush()
}
