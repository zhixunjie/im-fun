package websocket

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
)

// WriteMessage write a message（跟read反着来就好）
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.2
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.5
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.6
func (c *Conn) WriteMessage(msgType int, msg []byte) (err error) {
	// 1. write header
	err = c.WriteHeader(msgType, len(msg))
	if err != nil {
		return
	}
	// 2. write body
	err = c.WriteBody(msg)
	return
}

// WriteHeader write header frame.
func (c *Conn) WriteHeader(opcode int, length int) error {
	// 1. set first byte(8bit)
	firstByte := byte(0)
	firstByte = firstByte | fin          // fin
	firstByte = firstByte | byte(opcode) // opcode
	err := c.writer.WriteByte(firstByte)
	if err != nil {
		logrus.Errorf("err=%v", err)
		return err
	}

	// //////////////////////////////
	// 2. set second byte(8bit)
	writerBuff := make([]byte, 8) // TODO try to reduce gc
	secondByte := byte(0)
	switch {
	case length < 126:
		secondByte = secondByte | byte(length)
		// write second byte
		err = c.writer.WriteByte(secondByte)
		if err != nil {
			logrus.Errorf("err=%v", err)
			return err
		}
	case length < 65536:
		// write second byte
		secondByte = secondByte | 126
		err = c.writer.WriteByte(secondByte)
		if err != nil {
			logrus.Errorf("err=%v", err)
			return err
		}
		// write other byte
		binary.BigEndian.PutUint16(writerBuff[:2], uint16(length))
		_, err = c.writer.Write(writerBuff[:2])
		if err != nil {
			logrus.Errorf("err=%v", err)
			return err
		}
	default:
		// write second byte
		secondByte = secondByte | 127
		err = c.writer.WriteByte(secondByte)
		if err != nil {
			logrus.Errorf("err=%v", err)
			return err
		}

		// write other byte
		binary.BigEndian.PutUint64(writerBuff[:8], uint64(length))
		_, err = c.writer.Write(writerBuff[:8])
		if err != nil {
			logrus.Errorf("err=%v", err)
			return err
		}
	}
	return nil
}

// WriteBody write a message body.
func (c *Conn) WriteBody(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.writer.Write(b)
		if err != nil {
			logrus.Errorf("Write err=%v", err)
		}
	}
	return
}

// Flush flush writer buffer
func (c *Conn) Flush() error {
	return c.writer.Flush()
}
