package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"io"
	"strings"
)

var (
	keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	// ErrBadRequestMethod bad request method
	ErrBadRequestMethod = errors.New("bad method")
	// ErrNotWebSocket not websocket protocol
	ErrNotWebSocket = errors.New("not websocket protocol")
	// ErrBadWebSocketVersion bad websocket version
	ErrBadWebSocketVersion = errors.New("missing or bad WebSocket Version")
	// ErrChallengeResponse mismatch challenge response
	ErrChallengeResponse = errors.New("mismatch challenge/response")
)

// Upgrade Switching Protocols
// https://datatracker.ietf.org/doc/html/rfc6455#section-1.3
func Upgrade(closer io.ReadWriteCloser, rr *bufio.Reader, wr *bufio.Writer, req *Request) (conn *Conn, err error) {
	// check header
	err = verifyHeader(req)
	if err != nil {
		return nil, err
	}

	// upgrade
	challengeKey := req.Header.Get("Sec-Websocket-Key")
	_, _ = wr.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	_, _ = wr.WriteString("Upgrade: websocket\r\n")
	_, _ = wr.WriteString("Connection: Upgrade\r\n")
	_, _ = wr.WriteString("Sec-WebSocket-Accept: " + computeAcceptKey(challengeKey) + "\r\n\r\n")
	// _, _ = wr.WriteString("Sec-WebSocket-Protocol: " + computeAcceptKey(challengeKey) + "\r\n\r\n")
	if err = wr.Flush(); err != nil {
		logging.Errorf("Flush err=%v", err)
		return nil, err
	}

	return newConn(closer, rr, wr), nil
}

func computeAcceptKey(challengeKey string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(challengeKey))
	_, _ = h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func verifyHeader(req *Request) (err error) {
	if req.Method != "GET" {
		return ErrBadRequestMethod
	}
	if req.Header.Get("Sec-Websocket-Version") != "13" {
		return ErrBadWebSocketVersion
	}
	if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
		return ErrNotWebSocket
	}
	if !strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") {
		return ErrNotWebSocket
	}

	if req.Header.Get("Sec-Websocket-Key") == "" {
		return ErrChallengeResponse
	}

	return nil
}
