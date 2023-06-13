package bufio

import "io"

// 源码来源：Go 1.18.10，这里为其reader和write作能力加强

// Reader

// ResetBuffer reset reader & underlying byte array
func (b *Reader) ResetBuffer(r io.Reader, buf []byte) {
	b.reset(buf, r)
}

// ReadBytesN 保证从Reader中一次性读取N个字节的数据
func (b *Reader) ReadBytesN(buf []byte) error {
	_, err := io.ReadFull(b, buf)

	return err
}

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

// Writer

// ResetBuffer reset writer & underlying byte array
func (b *Writer) ResetBuffer(w io.Writer, buf []byte) {
	b.buf = buf
	b.err = nil
	b.n = 0
	b.wr = w
}

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
