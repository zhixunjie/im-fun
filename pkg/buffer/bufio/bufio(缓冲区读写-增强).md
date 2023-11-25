# 1. 源码阅读

> 先弄明白源码库的大概原理，然后再进行修改。

**bufio.Reader 和 bufio.Writer：**

- bufio.Reader，相当于为某个Reader(fd)在读取操作时，增加了缓冲区（存放一次性Read的内容），从而减少**系统调用Read**的调用次数。
- bufio.Writer，相当于为某个Writer(fd)在写入操作时，增加了缓冲区（存放即将要Write的内容），从而减少**系统调用Write**的调用次数。

## 1.1 Reader原理

> **数据流：bufio.Reader调用 Read => 底层Reader调用 Read => 保存到Reader的buf数组（预读到用户缓冲区） => 从buf中返回N个字节**

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

- **buf**：bufio.Reader 底层使用的实际内存。
- **rd**：bufio.Reader 底层使用使用的Reader（fd）。

----

2、**读取操作：** 参照[Reader.ReadByte函数](https://github.com/zhixunjie/im-fun/blob/584f7ec67b1140de3dcabc2bb6a73835421d0b9b/pkg/buffer/bufio/bufio.go#L258)的源码。

- 执行bufio.Reader的Read操作时，如果发现缓冲区buf的内容不足以满足本次读取的字节数，会调用fill函数，从而调用底层Reader(fd)的Read（系统调用）函数，一次性读取N个字节。
- 等到从底层Reader(fd)读取完毕后，再把执行最后的copy操作。

3、**读取过程中 r、w 值的变化过程**：执行读取操作，会增加r的值，执行写入操作，会增加w的值。

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

- 特点是：写入时，如果发现存放写入内容的缓冲区buf（用户缓冲区）空间不足时，会调用系统调用函数Write，把用户缓冲区的内容刷入（Flush）到内核缓存区后，从而释放buf（用户缓冲区）的内容以满足写入操作。

# 2. 功能加强

> 源码来源：基于Go 1.18.10，这里为其 reader 和 writer 做能力加强

**1、Reader.ReadBytesN()**：读取N个字节（bytes）。

- 虽然[Reader.read函数](https://github.com/zhixunjie/im-fun/blob/584f7ec67b1140de3dcabc2bb6a73835421d0b9b/pkg/buffer/bufio/bufio.go#L207)也能读取N个字节的数据，但是它不保证读取的字节数一定是等于N个字节，可能会小于N个字节。
- 因此，扩展了一个函数ReadBytesN()，能够保证从Reader中一次性读取N个字节的数据。

~~~go
func (b *Reader) ReadBytesN(buf []byte) error {
	_, err := io.ReadFull(b, buf)

	return err
}
~~~

**2、SetFdAndResetBuffer()：**设置**底层IO对象(fd)**，然后设置**底层IO对象**在读写操作（Read/Write系统调用）时用到的用户缓冲区。

- 优化：通过SetFdAndResetBuffer + Buffer Pool，能够在每条TCP链接中重复使用缓冲池子中的内存，从而使得长期运行的程序中大大减少GC的发生。
- 以下的代码，针对Reader和Writer都实现了SetFdAndResetBuffer操作。

~~~go
// SetFdAndResetBuffer reset reader & underlying byte array
func (b *Reader) SetFdAndResetBuffer(r io.Reader, buf []byte) {
	b.reset(buf, r)
}

// SetFdAndResetBuffer reset writer & underlying byte array
func (b *Writer) SetFdAndResetBuffer(w io.Writer, buf []byte) {
	b.buf = buf
	b.err = nil
	b.n = 0
	b.wr = w
}
~~~

**3、Reader.Pop()**：把Reader 的用户缓冲区的n个字节直接返回给用户。

- 如果用户调用Reader.Read()执行读取操作时，需要先make一个byte数组，然后传入Read()函数进行读取操作。
- 但是，如果用户调用Reader.Pop()执行读取操作时，就不再需要传入一个byte数组，直接复用Reader本身的缓冲区内存即可。

~~~go
// Pop 直接返回Reader的用户缓冲区的n个字节
// returns the next n bytes with advancing the reader. The bytes stop
// being valid at the next read call. If Pop returns fewer than n bytes, it
// also returns an error explaining why the read is short. The error is
// ErrBufferFull if n is larger than b's buffer size.
func (b *Reader) Pop(n int) ([]byte, error) {
	d, err := b.Peek(n)
	if err == nil {
		b.r += n
		return d, err
	}
	return nil, err
}
~~~

**4、Writer.Peek()**：原理跟Reader.Pop()一样。

~~~go
// Peek 直接返回Writer的用户缓冲区的n个字节
// returns the next n bytes with advancing the writer. The bytes stop
// being used at the next write call. If Peek returns fewer than n bytes, it
// also returns an error explaining why the read is short. The error is
// ErrBufferFull if n is larger than b's buffer size.
func (b *Writer) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrNegativeCount
	}
	if n > len(b.buf) {
		return nil, ErrBufferFull
	}
	for b.Available() < n && b.err == nil {
		b.Flush()
	}
	if b.err != nil {
		return nil, b.err
	}
	d := b.buf[b.n : b.n+n]
	b.n += n
	return d, nil
}
~~~

# 3. 汇总：项目用到的函数

reader：

~~~go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)    // 读取一行数据
func (b *Reader) ReadByte() (byte, error)                              // 读取一个字节
func (b *Reader) ReadBytesN(buf []byte)                                // 写入len(p)个字节
func (b *Reader) SetFdAndResetBuffer(r io.Reader, buf []byte)
func (b *Reader) Pop(n int) ([]byte, error) 
~~~

writer：

~~~go
func (b *Writer) WriteByte(c byte) error               // 写入一个字节
func (b *Writer) Write(p []byte) (nn int, err error)   // 写入len(p)个字节
func (b *Writer) Peek(n int) ([]byte, error)
func (b *Writer) SetFdAndResetBuffer(w io.Writer, buf []byte)
~~~

