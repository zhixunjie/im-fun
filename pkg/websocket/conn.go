package websocket

import (
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"io"
)

// Conn represents a WebSocket connection.
type Conn struct {
	closer  io.ReadWriteCloser
	reader  *bufio.Reader
	writer  *bufio.Writer
	maskKey []byte
}

func newConn(closer io.ReadWriteCloser, r *bufio.Reader, w *bufio.Writer) *Conn {
	return &Conn{closer: closer, reader: r, writer: w, maskKey: make([]byte, 4)}
}

func (c *Conn) Close() error {
	return c.closer.Close()
}
