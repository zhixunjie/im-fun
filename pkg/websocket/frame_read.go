package websocket

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// ReadMessage 读取整个消息
func (c *Conn) ReadMessage() (payload []byte, err error) {
	var isFin bool
	var finOp, n int
	var partPayload []byte
	var opCode int

	for {
		// read frame
		if isFin, opCode, partPayload, err = c.readFrame(); err != nil {
			return payload, err
		}
		// check opcode
		switch opCode {
		case BinaryMessage, TextMessage, continuationFrame:
			if isFin && len(payload) == 0 {
				return partPayload, nil
			}
			// continuation frame
			payload = append(payload, partPayload...)
			if opCode != continuationFrame {
				finOp = opCode
			}
			// final frame
			if isFin {
				opCode = finOp
				return payload, err
			}
		case PingMessage: // PING
			if err = c.WriteMessage(PongMessage, partPayload); err != nil {
				return
			}
		case PongMessage: // PONG
		case CloseMessage: // CLOSE
			err = ErrMessageClose
			return payload, err
		default:
			err = fmt.Errorf("unknown frame")
			logging.Errorf("err=%v,isFin=%t,opCode=%d", err, isFin, opCode)
			return payload, err
		}
		if n > continuationFrameMaxRead {
			err = ErrMessageMaxRead
			return payload, err
		}
		n++
	}
}

// 读取帧
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.2
func (c *Conn) readFrame() (isFin bool, opCode int, payload []byte, err error) {
	// 1. get first byte(8bit)
	var firstByte byte
	firstByte, err = c.reader.ReadByte()
	if err != nil {
		return
	}

	// 2. read bit: fin
	fin := firstByte & fin
	isFin = fin != 0 // equal 1 means final frag

	// 3. read bit：rsv
	// rsv MUST be 0
	if rsv := firstByte & (rsv1 | rsv2 | rsv3); rsv != 0 {
		err = errors.New("rsv not allow")
		return false, 0, nil, err
	}
	// 4. read bit: opcode
	opCode = int(firstByte & opcode)

	// //////////////////////////
	// 2. get second byte(8bit)
	var secondByte byte
	secondByte, err = c.reader.ReadByte()
	if err != nil {
		return
	}

	// 2. read bit: mark
	// Defines whether the "Payload data" is masked.
	masked := (secondByte & mark) != 0

	// 3. read bit: payload's length
	// note: payload's length maybe 1 byte，maybe multiple bytes
	// - therefore, we read 1 byte at first
	// - if we find that value exceed 125，length need more bytes to save
	// - we read the next 2 bytes, or the next 8 bytes
	var payloadLen int64
	lenVal := int64(secondByte & payloadLength)

	// 4. get payload length（）
	var readerBuffer []byte
	switch {
	case lenVal < 126:
		// 1) if 0-125, that is the payload length
		payloadLen = lenVal
	case lenVal == 126:
		// 2) If 126, the following 2 bytes interpreted as a 16-bit unsigned integer are the payload length.
		if readerBuffer, err = c.reader.Pop(2); err != nil {
			logging.Errorf("Pop(2) readerBuffer err=%v", err)
			return
		}
		payloadLen = int64(binary.BigEndian.Uint16(readerBuffer))
	case lenVal == 127:
		// 3)  If 127, the following 8 bytes interpreted as a 64-bit unsigned integer
		// (the most significant bit MUST be 0) are the payload length.
		if readerBuffer, err = c.reader.Pop(8); err != nil {
			logging.Errorf("Pop(8) readerBuffer err=%v", err)
			return
		}
		payloadLen = int64(binary.BigEndian.Uint16(readerBuffer))
	}
	if payloadLen < 0 {
		logging.Errorf("payloadLen not allow")
		return
	}

	// 3. read bit: masked key
	// All frames sent from the client to the server are masked by a 32-bit value that is contained within the frame.
	// masked key is used to masked payload
	var maskKey []byte
	if masked {
		if maskKey, err = c.reader.Pop(4); err != nil {
			logging.Errorf("Pop(4) maskKey err=%v", err)
			return
		}
		if c.maskKey == nil {
			c.maskKey = make([]byte, 4)
		}
		copy(c.maskKey, maskKey)
	}
	// 4. read payload（finally，OMG）
	// https://datatracker.ietf.org/doc/html/rfc6455#section-5.3
	if payloadLen > 0 {
		if payload, err = c.reader.Pop(int(payloadLen)); err != nil {
			logging.Errorf("Pop(%v) payload err=%v", payloadLen, err)
			return
		}
		if masked {
			maskBytes(c.maskKey, 0, payload)
		}
	}
	return
}

func maskBytes(key []byte, pos int, payload []byte) int {
	for i := range payload {
		payload[i] ^= key[pos&3]
		pos++
	}
	return pos & 3
}
