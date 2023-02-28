package websocket

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"golang.org/x/net/websocket"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	data := []byte("hello client")
	var err error

	// listen
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.FailNow()
	}
	for {
		conn, err := listen.Accept()
		assert.Nil(t, err)
		rd := bufio.NewReader(conn)
		wr := bufio.NewWriter(conn)
		req, err := ReadRequest(rd)
		assert.Nil(t, err)
		assert.Equal(t, "/sub", req.RequestURI)

		// upgrade
		ws, err := Upgrade(conn, rd, wr, req)
		assert.Nil(t, err)

		for {
			err = ws.WriteMessage(TextMessage, data)
			assert.Nil(t, err)

			err = ws.Flush()
			assert.Nil(t, err)

			buf, err := ws.ReadMessage()
			assert.Nil(t, err)
			fmt.Printf("%s\n", buf)

			time.Sleep(2 * time.Second)
		}
	}
}

func TestClient(t *testing.T) {
	ws, err := websocket.Dial("ws://127.0.0.1:8080/sub", "", "*")
	if err != nil {
		t.FailNow()
	}

	for {
		// receive binary frame
		var buf []byte
		err = websocket.Message.Receive(ws, &buf)
		assert.Nil(t, err)
		fmt.Printf("%s\n", buf)

		// send binary frame
		data := []byte("hello server")
		err = websocket.Message.Send(ws, data)
		assert.Nil(t, err)
	}
}
