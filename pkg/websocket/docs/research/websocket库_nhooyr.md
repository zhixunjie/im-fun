> 记录用于实现参考的关键方法

- https://github.com/nhooyr/websocket
- https://pkg.go.dev/nhooyr.io/websocket

upgrade：

~~~go
func accept(w http.ResponseWriter, r *http.Request, opts *AcceptOptions)
~~~

read：

~~~go
func readFrameHeader(r *bufio.Reader, readBuf []byte) (h header, err error) 
n, err := mr.c.readFramePayload(mr.ctx, p)
~~~

write：

~~~~go
func writeFrameHeader(h header, w *bufio.Writer, buf []byte) (err error) 
func (c *Conn) writeFramePayload(p []byte) (n int, err error)
~~~~

