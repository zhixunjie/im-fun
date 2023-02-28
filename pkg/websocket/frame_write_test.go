package websocket

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"os"
	"testing"
)

func TestFrameWrite(t *testing.T) {
	filePath := "a.txt"
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(fd)
	wr := bufio.NewWriter(fd)

	// write
	conn := newConn(nil, rd, wr)
	data := []byte("Hello Client")
	err = conn.WriteMessage(BinaryMessage, data)
	if err != nil {
		fmt.Println("WriteMessage", err)
		return
	}
	fmt.Println(conn.Flush())
}

func TestFrameRead(t *testing.T) {
	filePath := "a.txt"
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(fd)
	wr := bufio.NewWriter(fd)

	// read
	conn := newConn(nil, rd, wr)
	buf, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("ReadMessage", err)
		return
	}
	fmt.Printf("%s\n", buf)
}
