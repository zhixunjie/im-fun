package operation

import (
	"bufio"
	"context"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net"
	"sync/atomic"
)

func Start(userId int64, addr string) {
	//time.Sleep(time.Duration(rand.Intn(120)) * time.Second)

	// dial to server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logging.Errorf("net.Dial(%s) error(%v)", addr, err)
		return
	}

	// quit
	atomic.AddInt64(&aCount, 1)
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
		atomic.AddInt64(&aCount, -1)
	}()

	// init
	ctx := context.Background()
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	// auth
	seq := int32(0)
	if err = Auth(rd, wr, userId); err != nil {
		return
	}
	seq++

	// writer
	go Writer(ctx, &seq, wr, userId, quit)

	// reader
	_ = Reader(ctx, conn, rd, userId, quit)
}
