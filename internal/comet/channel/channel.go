package channel

import (
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	"github.com/zhixunjie/im-fun/pkg/logging"
	newtimer "github.com/zhixunjie/im-fun/pkg/time"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"net"
	"time"
)

// Channel 每一个Channel代表一个长连接
type Channel struct {
	ConnComponent

	// the Room that Channel(User) is in
	Room   *Room
	signal chan *protocol.Proto
	Next   *Channel
	Prev   *Channel

	// use bufio to reuse buffer
	Writer         *bufio.Writer
	Reader         *bufio.Reader
	ConnReadWriter protocol.ConnReadWriter

	// user info
	ClientIp string
	UserInfo *pb.TcpUserInfo

	// use to get proto，reduce GC
	ProtoAllocator Ring
}

// NewChannel new a channel.
func NewChannel(conf *conf.Config, conn *net.TCPConn, traceId int64, connType ConnType, readerPool, writerPool *bytes.Pool, timerPool *newtimer.Timer) *Channel {
	// init channel
	ch := &Channel{
		// set ConnComponent
		ConnComponent: ConnComponent{
			TraceId:    traceId,
			ConnType:   connType,
			Conn:       conn,
			WsConn:     nil,
			WriterPool: writerPool,
			ReaderPool: readerPool,
			writeBuf:   writerPool.Get(),
			readBuf:    readerPool.Get(),
			TimerPool:  timerPool,
			Trd:        nil,
		},
		signal: make(chan *protocol.Proto, conf.Protocol.Proto.ChannelSize),
	}

	// set ProtoAllocator
	ch.ProtoAllocator.Init(uint64(conf.Protocol.Proto.AllocatorSize))

	// set user ip
	ch.ClientIp, _, _ = net.SplitHostPort(conn.RemoteAddr().String())

	// set buffer
	// 底层执行IO操作的是conn，缓冲区为ch.readBuf.Bytes()
	reader := new(bufio.Reader)
	reader.SetFdAndResetBuffer(conn, ch.readBuf.Bytes())
	writer := new(bufio.Writer)
	writer.SetFdAndResetBuffer(conn, ch.writeBuf.Bytes())
	// set connection's reader and writer
	ch.Reader = reader
	ch.Writer = writer
	ch.ConnReadWriter = protocol.NewTcpConnReaderWriter(reader, writer)

	return ch
}

func (c *Channel) SetWebSocketConnReaderWriter(wsConn *websocket.Conn) {
	c.ConnComponent.WsConn = wsConn
	c.ConnReadWriter = protocol.NewWsConnReaderWriter(wsConn)
}

func (c *Channel) CleanPath1() {
	logHead := fmt.Sprintf("[traceId=%v] CleanPath1|", c.TraceId)

	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
		logging.Info(logHead + "WsConn.Close")
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
			logging.Info(logHead + "Conn.Close")
		}
	}
	// 2. 把buffer重新放回到Pool(read & write)
	if c.writeBuf != nil {
		c.WriterPool.Put(c.writeBuf)
		logging.Info(logHead + "WriterPool.Put")
	}
	if c.readBuf != nil {
		c.ReaderPool.Put(c.readBuf)
		logging.Info(logHead + "ReaderPool.Put")
	}
	// 3. 把timer从Pool里面删除
	if c.Trd != nil {
		c.TimerPool.Del(c.Trd)
	}
}

// CleanPath2 Read协程结束时，需要执行的清理动作
func (c *Channel) CleanPath2() {
	logHead := fmt.Sprintf("[traceId=%v] CleanPath2|", c.TraceId)

	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
		logging.Info(logHead + "WsConn.Close")
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
			logging.Info(logHead + "Conn.Close")
		}
	}

	// 2. 把buffer重新放回到Pool(only read buffer)
	// note: writePool's buffer will be released in Server.dispatch
	if c.readBuf != nil {
		c.ReaderPool.Put(c.readBuf)
		logging.Info(logHead + "ReaderPool.Put")
	}

	// 3. 把timer从Pool里面删除
	if c.Trd != nil {
		c.TimerPool.Del(c.Trd)
	}

	// 4. SendFinish
	if c.signal != nil {
		c.SendFinish(logHead)
		logging.Info(logHead + "SendFinish")
	}
}

// CleanPath3 Write协程结束时，需要执行的清理动作
func (c *Channel) CleanPath3() {
	logHead := fmt.Sprintf("[traceId=%v] CleanPath3|", c.TraceId)

	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
		logging.Info(logHead + "WsConn.Close")
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
			logging.Info(logHead + "Conn.Close")
		}
	}
	// 2. 把buffer重新放回到Pool(only write buffer)
	if c.writeBuf != nil {
		c.WriterPool.Put(c.writeBuf)
		logging.Info(logHead + "WriterPool.Put")
	}
}

// ConnComponent 每一条连接需要用到的组件
type ConnComponent struct {
	TraceId  int64
	ConnType ConnType

	// Connection(fd)
	Conn   *net.TCPConn
	WsConn *websocket.Conn

	// 分配buffer池子
	WriterPool *bytes.Pool
	ReaderPool *bytes.Pool
	// 从池子分配得到Buffer
	writeBuf *bytes.Buffer
	readBuf  *bytes.Buffer

	// 分配定时器的池子
	TimerPool *newtimer.Timer
	// 从池子分配得到Timer
	Trd *newtimer.TimerData

	// 心跳相关
	LastHb     time.Time     // 上一次接收到心跳的时间
	HbExpire   time.Duration // 心跳超时的时间
	HbInterval time.Duration // 心跳续约频率
}
