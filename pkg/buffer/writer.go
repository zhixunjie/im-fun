package buffer

// Writer writer.
type Writer struct {
	n   int // n means buffer has been written
	buf []byte
}

// NewWriterSize new a writer with size.
func NewWriterSize(n int) *Writer {
	return &Writer{buf: make([]byte, n)}
}

// Len buffer has been written
func (w *Writer) Len() int {
	return w.n
}

// Size buffer
func (w *Writer) Size() int {
	return len(w.buf)
}

// Reset buffer
func (w *Writer) Reset() {
	w.n = 0
}

// Buffer return buffer's []byte
// 返回已写入的缓冲区数据
func (w *Writer) Buffer() []byte {
	return w.buf[:w.n]
}

// Peek buffer
// 预留缓冲区的一段内存（即将会被写入数据的内存）
func (w *Writer) Peek(n int) []byte {
	var buf []byte
	w.grow(n)
	buf = w.buf[w.n : w.n+n]
	w.n += n
	return buf
}

// Write buffer
func (w *Writer) Write(p []byte) {
	w.grow(len(p))
	w.n += copy(w.buf[w.n:], p)
}

func (w *Writer) grow(n int) {
	var buf []byte
	if w.n+n < len(w.buf) {
		return
	}
	buf = make([]byte, 2*len(w.buf)+n)
	copy(buf, w.buf[:w.n])
	w.buf = buf
}
