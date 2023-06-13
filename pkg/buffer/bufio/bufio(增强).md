# 1. 源码阅读

> 先弄明白源码库的大概原理，然后再进行修改。

**bufio.Reader 和 bufio.Writer：**

- bufio.Reader，相当于为某个Reader(fd)在读取操作时，增加了缓冲区（存放一次性Read的内容），从而减少**系统调用Read**的调用次数。
- bufio.Writer，相当于为某个Writer(fd)在写入操作时，增加了缓冲区（存放即将要Write的内容），从而减少**系统调用Write**的调用次数。

## 1.1 Reader原理

> bufio.Reader调用 Read => 底层Reader调用 Read => 保存到buf数组 => 从buf中返回N个字节

~~~go
// Reader implements buffering for an io.Reader object.
type Reader struct {
	buf          []byte
	rd           io.Reader // reader provided by the client
	r, w         int       // buf read and write positions
	err          error
	lastByte     int // last byte read for UnreadByte; -1 means invalid
	lastRuneSize int // size of last rune read for UnreadRune; -1 means invalid
}
~~~

1、**字段说明：**

- **buf**：bufio.Reader底层使用的实际内存。
- **rd**：bufio.Reader底层使用使用的Reader（fd）。

----

2、**读取操作：** 参照[Reader.ReadByte函数](https://github.com/zhixunjie/im-fun/blob/584f7ec67b1140de3dcabc2bb6a73835421d0b9b/pkg/buffer/bufio/bufio.go#L258)的源码。

- 执行bufio.Reader的Read操作时，如果发现缓冲区buf的内容不足以满足本次读取的字节数，会调用fill函数，从而调用底层Reader(fd)的Read（系统调用）函数，一次性读取N个字节。
- 等到从底层Reader(fd)读取完毕后，再把执行最后的copy操作。

3、**r、w值的变化过程**：执行读取操作，会增加r的值，执行写入操作，会增加w的值。

- **w的值比r的值大**：说明还能从buf中读取操作。
- **w的值等于r的值**：
  - 说明buf的内容已经被全部读取完毕，需要调用系统调用Read函数，从fd中读取内容。
  - 等到从fd读取完毕后，就会增加w的值，此时w的值恢复为大于r的值。

## 1.2 Writer原理

~~~go
// Writer implements buffering for an io.Writer object.
// If an error occurs writing to a Writer, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
type Writer struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}
~~~

大致跟Reader相同，参照[Writer.WriteByte函数](https://github.com/zhixunjie/im-fun/blob/584f7ec67b1140de3dcabc2bb6a73835421d0b9b/pkg/buffer/bufio/bufio.go#L687)的源码。

- 特点是：写入时，如果发现存放写入内容的缓冲区buf（用户缓冲区）空间不足时，会把调用系统调用函数Write，把当前缓冲区的刷入（Flush）到内核缓存区后，从而释放buf（用户缓冲区）的内容以满足写入操作。

# 2. 功能加强

> 源码来源：Go 1.18.10，这里为其reader和writer作能力加强

1、读取N个字节（bytes）：**read N bytes**

- 虽然[Reader.read函数](https://github.com/zhixunjie/im-fun/blob/584f7ec67b1140de3dcabc2bb6a73835421d0b9b/pkg/buffer/bufio/bufio.go#L207)也能读取N个字节的数据，但是它不保证读取的字节数一定是等于N个字节，可能会小于N个字节。
- 因此，扩展了一个函数ReadBytesN()，能够保证从Reader中一次性读取N个字节的数据。

~~~go
func (b *Reader) ReadBytesN(buf []byte) error {
	_, err := io.ReadFull(b, buf)

	return err
}
~~~

2、重置底层Reader对象，以及对底层对象读写操作（Read/Write系统调用）时用到的用户缓冲区。

- 优化：通过Reset + Buffer Pool，能够在长期运行的程序中大大减少GC的发生。
- 以下的代码，针对Reader和Writer都实现了Reset操作。

~~~go
func (b *Reader) ResetBuffer(r io.Reader, buf []byte) {
	b.reset(buf, r)
}

func (b *Writer) ResetBuffer(w io.Writer, buf []byte) {
	b.buf = buf
	b.err = nil
	b.n = 0
	b.wr = w
}
~~~

# 3. 汇总：项目用到的函数

reader：

~~~go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)    // 读取一行数据
func (b *Reader) ReadByte() (byte, error)                              // 读取一个字节
func (b *Reader) ReadBytesN(buf []byte)                                // 写入len(p)个字节
~~~

writer：

~~~go
func (b *Writer) WriteByte(c byte) error               // 写入一个字节
func (b *Writer) Write(p []byte) (nn int, err error)   // 写入len(p)个字节
~~~

