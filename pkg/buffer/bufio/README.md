a little bit change for bufio struct in Go1.18

note：it can only invoke the function that will grow the underlying buffer

---

read function that can be invoked：

~~~go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
~~~

write function that can be invoked：

~~~go
~~~

other funcint that can be invoked：

~~~go
ResetBuffer(r io.Reader, buf []byte)
~~~

