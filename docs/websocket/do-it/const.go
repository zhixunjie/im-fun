package websocket

import "errors"

const (
	// Frame Protocol, Section 5.2 of RFC 6455
	// https://datatracker.ietf.org/doc/html/rfc6455#section-5.2

	// first byte
	fin    = 1 << 7
	rsv1   = 1 << 6
	rsv2   = 1 << 5
	rsv3   = 1 << 4
	opcode = 0x0f

	// second bit
	mark          = 1 << 7
	payloadLength = 0x7f

	continuationFrame        = 0
	continuationFrameMaxRead = 100
)

// The message types are defined in RFC 6455, section 11.8.
// https://datatracker.ietf.org/doc/html/rfc6455#section-11.8
// https://datatracker.ietf.org/doc/html/rfc6455#section-5.5
// note：
// 1) 3-7 are reserved for further non-control frames.
// 2) 11-16 are reserved for further control frames.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

var (
	// ErrMessageClose close control message
	ErrMessageClose = errors.New("close control message")
	// ErrMessageMaxRead continuation frame max read
	ErrMessageMaxRead = errors.New("continuation frame max read")
)
