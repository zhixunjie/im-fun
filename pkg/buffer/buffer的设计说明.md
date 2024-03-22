# Buffer Pool

~~~shell
# 关系图
Pool -> Buffer
     		-> Get -> get a Buffer from the Pool
     		-> Put -> put back a Buffer to the Pool
~~~

~~~go
// 内存池
type Pool struct {
	lock     sync.Mutex
	free     *Buffer // check this detail in batchNew
	bufSize  int     // each Buffer Size
	batchNum int     // create Buffer continuously
}

// 一个内存块
type Buffer struct {
	buf  []byte
	next *Buffer // Point to the next free Buffer
}
~~~

具体说明：

**1、Pool & Buffer：** 内存池，用于分配内存。

- Pool 负责管理和分配 Buffer（池子里面的一个「内存单元」，就是一个「Buffer对象」）。
- **Pool 是如何分配内存的？ **  How to allocate memory？
  - Pool 的本质是通过链表方式把「内存单元」连接在一起。
  - 每次分配内存时，从链表头部取出一个Buffer对象即可。
  - 如果发现Buffer Pool内没有Buffer，需要预先分配一大段内存再进行切分（批量创建buffer）。
    - 相对于golang自带的sync.Pool， 好处就是批量New，而不是一个个去New。

---

**2、Pool Hash：**  池子分片，基于「哈希取余」的方式进行池子分配。好处：减少单个池子的 Mutex 冲突。

~~~go
// 利用Hash算法，均摊池子的请求流量
type Hash struct {
	Readers []Pool
	Writers []Pool
	options *Options
}
~~~

# Bufio

> 核心点：减少系统调用次数、减少磁盘操作次数。

Bufio：为某个 fd 添加用户缓冲区的读写操作，[改造版的Bufio](./bufio/bufio(缓冲区读写-增强).md)。

- **Bufio 的内存复用了 [Buffer Pool](#Buffer Pool)的内存，从而减少每个TCP的IO读写带来的Buffer GC。**
    - 由于每个TCP连接（conn fd）都需要附带上 Bufio 的用户缓冲区，频繁进行内存的创建和销毁，对于申请内存和GC都是要消耗性能的。
    - 所以，用户缓冲区的内存交由Buffer Pool去管理。
- Bufio 为底层的 read 和 write 操作附上用户缓冲区（从而减少系统调用 read/write 的次数）。
    - 相当于让TCP Reader(conn)的读写带上了用户缓冲区（相当于C语言的标准IO函数的用户缓冲区），从而减少conn的系统调用 read/write 的次数。

> **备注：如何减少磁盘操作次数？** 指定Socket的读写缓冲区大小，当缓冲区满后才会真正执行磁盘的操作。
>
> - **具体见：延迟写.md**  
> - SetReadBuffer：sets the size of the operating system's receive buffer associated with the connection.
> - SetWriteBuffer：sets the size of the operating system's transmit buffer associated with the connection.

~~~go
if err = conn.SetReadBuffer(server.conf.Connect.TCP.Rcvbuf); err != nil {
  logging.Errorf(logHead+"conn.SetReadBuffer() error=%v", err)
  return
}

if err = conn.SetWriteBuffer(server.conf.Connect.TCP.Sndbuf); err != nil {
  logging.Errorf(logHead+"conn.SetWriteBuffer() error=%v", err)
  return
}
~~~

