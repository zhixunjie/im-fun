package buffer

// BufferWriter：申请一大段内存，然后向内存执行写入/重置/Peek等操作
// - 每次写入后，都会增加n的值。如果写入时，发现buffer的空间不足，则调用grow自动对buffer进行扩容。

type BufferWriter struct {
	n   int // n means buffer has been written
	buf []byte
}

func NewWriterSize(n int) *BufferWriter {
	return &BufferWriter{buf: make([]byte, n)}
}

// Len 获取buffer已经被写入的大小
func (w *BufferWriter) Len() int {
	return w.n
}

// Size 获取buffer的大小
func (w *BufferWriter) Size() int {
	return len(w.buf)
}

// Reset 重置buffer
func (w *BufferWriter) Reset() {
	w.n = 0
}

// Buffer 获取buffer已写入的全部护具
func (w *BufferWriter) Buffer() []byte {
	return w.buf[:w.n]
}

// Peek 获取缓冲区的一块内存地址（即：返回的内存会被调用方执行写入操作）
func (w *BufferWriter) Peek(n int) []byte {
	var buf []byte
	w.grow(n)
	buf = w.buf[w.n : w.n+n]
	w.n += n
	return buf
}

// Write 把数据copy到Writer的缓冲区
func (w *BufferWriter) Write(p []byte) {
	w.grow(len(p))
	w.n += copy(w.buf[w.n:], p)
}

// 内存自增长：自动扩展内存
func (w *BufferWriter) grow(n int) {
	var buf []byte
	if w.n+n < len(w.buf) {
		return
	}
	buf = make([]byte, 2*len(w.buf)+n)
	copy(buf, w.buf[:w.n])
	w.buf = buf
}
