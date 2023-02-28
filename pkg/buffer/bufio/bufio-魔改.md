a little bit change for bufio struct in Go1.18

Why？add some function to it for promoting performance，such as：

**read N bytes：**bufio just only offer the ReadByte() function 

~~~go
func (b *Reader) ReadBytesN(buf []byte) error {
	_, err := io.ReadFull(b, buf)

	return err
}
~~~

**reuse underlying buffer：**through these method and buffer pool，we can reduce gc in  a long term running.

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

---

**note：**it can only invoke the function that will grow the underlying buffer

read function that can be invoked：

~~~go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
func (b *Reader) ReadByte() (byte, error) 
func (b *Reader) ReadBytesN(buf []byte) 
~~~

write function that can be invoked：

~~~go
func (b *Writer) WriteByte(c byte) 
func (b *Writer) Write(p []byte) 
~~~

other funcint that can be invoked：

~~~go
ResetBuffer(r io.Reader, buf []byte)
~~~

