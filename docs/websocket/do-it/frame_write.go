package websocket

import (
	"encoding/binary"
)

// WriteMessage write a message by type.
func (c *Conn) WriteMessage(msgType int, msg []byte) (err error) {
	if err = c.WriteHeader(msgType, len(msg)); err != nil {
		return
	}
	err = c.WriteBody(msg)
	return
}

// WriteHeader write header frame.
func (c *Conn) WriteHeader(msgType int, length int) (err error) {
	var h []byte
	if h, err = c.writer.Peek(2); err != nil {
		return
	}
	// 1.First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	h[0] = 0
	h[0] = h[0] | (fin | byte(msgType))
	// 2.Second byte. Mask/Payload len(7bits)
	h[1] = 0
	switch {
	case length <= 125:
		// 7 bits
		h[1] |= byte(length)
	case length < 65536:
		// 16 bits
		h[1] |= 126
		if h, err = c.writer.Peek(2); err != nil {
			return
		}
		binary.BigEndian.PutUint16(h, uint16(length))
	default:
		// 64 bits
		h[1] |= 127
		if h, err = c.writer.Peek(8); err != nil {
			return
		}
		binary.BigEndian.PutUint64(h, uint64(length))
	}
	return
}

// WriteBody write a message body.
func (c *Conn) WriteBody(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.writer.Write(b)
	}
	return
}

// Peek write peek.
func (c *Conn) Peek(n int) ([]byte, error) {
	return c.writer.Peek(n)
}

// Flush flush writer buffer
func (c *Conn) Flush() error {
	return c.writer.Flush()
}
