package channel

import (
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/buffer/bufio"
	"github.com/zhixunjie/im-fun/pkg/buffer/bytes"
	newtimer "github.com/zhixunjie/im-fun/pkg/time"
	"github.com/zhixunjie/im-fun/pkg/websocket"
	"net"
)

// Channel 每一个Channel代表一个长连接
type Channel struct {
	ConnComponent

	// The Room that Channel(User) is In
	Room   *Room
	signal chan *protocol.Proto
	Next   *Channel
	Prev   *Channel

	// use bufio to reuse buffer
	Writer           *bufio.Writer
	Reader           *bufio.Reader
	ConnReaderWriter protocol.ConnReaderWriter

	// User Info
	UserInfo *UserInfo

	// use to get proto，reduce GC
	ProtoAllocator Ring
}

// NewChannel new a channel.
func NewChannel(conf *conf.Config, conn *net.TCPConn, connectionType int, readerPool, writerPool *bytes.Pool, timerPool *newtimer.Timer) *Channel {
	ch := new(Channel)
	ch.ProtoAllocator.Init(uint64(conf.Protocol.ClientProtoNum))
	ch.UserInfo = new(UserInfo)
	ch.signal = make(chan *protocol.Proto, conf.Protocol.ServerProtoNum)
	ch.Reader = new(bufio.Reader)
	ch.Writer = new(bufio.Writer)
	ch.ConnReaderWriter = protocol.NewTcpConnReaderWriter(ch.Reader, ch.Writer)

	// set ConnComponent
	ch.ConnComponent.ConnType = connectionType
	ch.ConnComponent.ReaderPool = readerPool
	ch.ConnComponent.WriterPool = writerPool
	ch.ConnComponent.TimerPool = timerPool
	ch.ConnComponent.Conn = conn

	// set reader's buffer
	var rb = readerPool.Get()
	ch.Reader.SetFdAndResetBuffer(conn, rb.Bytes())
	// set writer's buffer
	var wb = writerPool.Get()
	ch.Writer.SetFdAndResetBuffer(conn, wb.Bytes())
	// set user ip
	ch.UserInfo.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())

	return ch
}

func (c *Channel) CleanPath1() {
	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
	}
	// 2. 把buffer重新放回到Pool(read & write)
	if c.writeBuf != nil {
		c.WriterPool.Put(c.writeBuf)
	}
	if c.readBuf != nil {
		c.ReaderPool.Put(c.readBuf)
	}
	// 3. 把timer从Pool里面删除
	if c.Trd != nil {
		c.TimerPool.Del(c.Trd)
	}
}

func (c *Channel) CleanPath2() {
	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
	}
	// 2. 把buffer重新放回到Pool(only read buffer)
	// writePool's buffer will be released  in Server.dispatchTCP()
	if c.readBuf != nil {
		c.ReaderPool.Put(c.readBuf)
	}
	// 3. 把timer从Pool里面删除
	if c.Trd != nil {
		c.TimerPool.Del(c.Trd)
	}
	// 4. SendFinish
	if c.signal != nil {
		c.SendFinish()
	}
}

func (c *Channel) CleanPath3() {
	// 1. 关闭连接（一旦关闭连接，读写操作都会失败）
	if c.WsConn != nil {
		_ = c.WsConn.Close() // 它会把原本的Conn也给关闭掉
	} else {
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
	}
	// 2. 把buffer重新放回到Pool(only write buffer)
	if c.writeBuf != nil {
		c.ReaderPool.Put(c.writeBuf)
	}
}

const (
	ConnectionTypeTcp = iota + 1
	ConnectionTypeWebSocket
)

// ConnComponent 每一条连接需要用到的组件
type ConnComponent struct {
	ConnType int

	// Connection(fd)
	Conn   *net.TCPConn
	WsConn *websocket.Conn

	// 分配buffer池子
	WriterPool *bytes.Pool
	ReaderPool *bytes.Pool
	// 分配成功：得到Buffer
	writeBuf *bytes.Buffer
	readBuf  *bytes.Buffer

	// 分配定时器的池子
	TimerPool *newtimer.Timer
	// 分配成功：得到Timer
	Trd *newtimer.TimerData
}
