package bytes

import (
	"bytes"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"testing"
)

func TestPoolHash(t *testing.T) {
	poolHash := NewHash(&Options{
		ReadPool: PoolOptions{
			PoolNum:  1,
			BatchNum: 2,
			BufSize:  20,
		},
		WritePool: PoolOptions{
			PoolNum:  1,
			BatchNum: 2,
			BufSize:  20,
		},
	})

	pool := poolHash.ReaderPool(101)
	buf := pool.Get()
	byteArr := buf.Bytes()
	// fmt.Printf("%s\n", byteArr)

	// use Bufio to read connection
	// connReader := bytes.NewReader([]byte("abc\ndef\n")) // suppose this is a connection fd
	connReader := bytes.NewReader([]byte("abc\ndef")) // suppose this is a connection fd
	bufioReader := bufio.Reader{}
	bufioReader.SetFdAndResetBuffer(connReader, byteArr)

	// Read
	// 当connReader中的字符串最后没有换行，下面的输出会有bug，这个属于bufio原有的问题
	line1, isPrefix, err := bufioReader.ReadLine()
	fmt.Printf("line=%s,isPrefix=%v,err=%v\n", line1, isPrefix, err)
	line2, isPrefix, err := bufioReader.ReadLine()
	fmt.Printf("line=%s,isPrefix=%v,err=%v\n", line1, isPrefix, err)
	// print line
	fmt.Printf("%s\n", line1)
	fmt.Printf("%s\n", line2)
}
