package bufio

import "io"

// 源码来源：Go 1.18.10，这里为其reader和write作能力加强

// Reader

// ResetBuffer reset reader & underlying byte array
func (b *Reader) ResetBuffer(r io.Reader, buf []byte) {
	b.reset(buf, r)
}

func (b *Reader) ReadBytesN(buf []byte) error {
	_, err := io.ReadFull(b, buf)

	return err
}

// Pop 直接返回Reader的缓冲区（返回n个字节）
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

// TODO 待实现
// Peak 直接返回Writer的缓冲区（返回n个字节）
// func (b *WriterPool) Peak(n int) ([]byte, error) {
// 	if n < 0 {
// 		return nil, ErrNegativeCount
// 	}
// 	if n > len(b.buf) {
// 		return nil, ErrBufferFull
// 	}
// 	for b.Available() < n && b.err == nil {
// 		b.flush()
// 	}
// 	if b.err != nil {
// 		return nil, b.err
// 	}
// 	d := b.buf[b.n : b.n+n]
// 	b.n += n
// 	return d, nil
// }
