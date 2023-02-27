package buffer

import (
	"bytes"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"testing"
)

func TestPoolHash(t *testing.T) {
	poolHash := NewPoolHash(Options{
		r: PoolOptions{
			PoolNum:  1,
			BatchNum: 2,
			BufSize:  20,
		},
		w: PoolOptions{
			PoolNum:  1,
			BatchNum: 2,
			BufSize:  20,
		},
	})

	pool := poolHash.Reader(101)
	buf := pool.Get()
	byteArr := buf.Bytes()
	fmt.Printf("%s\n", byteArr)

	// Bufio
	bufioReader := bufio.Reader{}
	connReader := bytes.NewReader([]byte("abc\ndef")) // suppose this is a connection fd
	bufioReader.ResetBuffer(connReader, byteArr)

	// Read
	line1, _, _ := bufioReader.ReadLine()
	fmt.Printf("%s\n", line1)
	line2, _, _ := bufioReader.ReadLine()
	fmt.Printf("%s\n", line2)
}
