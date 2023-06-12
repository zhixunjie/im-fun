package buffer

// Buffer：缓冲区，参考：bytes.Buffer
// - 每个缓冲区代表一段指定大小的内存
// - 缓冲区与缓冲区之间，使用链表连接在一起

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
