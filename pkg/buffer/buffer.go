package buffer

type Buffer struct {
	buf  []byte
	next *Buffer // Point to the next free Buffer
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		buf:  make([]byte, size),
		next: nil,
	}
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}
