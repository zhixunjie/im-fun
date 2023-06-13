# 参考

- https://github.com/halfrost/Halfrost-Field/blob/master/contents/Protocol/WebSocket.md
- https://en.wikipedia.org/wiki/WebSocket

# overview

The handshake from the client looks as follows:

~~~
GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13
~~~

 The handshake from the server looks as follows:

~~~
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat
~~~

# 1. 客户端

- PHP：https://www.kancloud.cn/zhixunjie/swoole_websocket/363058

# 2. 服务端

第三方API：

- https://pkg.go.dev/nhooyr.io/websocket

RFC：

- https://datatracker.ietf.org/doc/html/rfc6455
- https://datatracker.ietf.org/doc/html/rfc6455#section-5.2 frame
- https://datatracker.ietf.org/doc/html/rfc6455#section-1.3 upgrade

**我的实现：**

- https://github.com/zhixunjie/im-fun/tree/master/pkg/websocket

**第三方库：**可以用于参考

- [nhooyr.md](research/nhooyr.md)
- https://github.com/gorilla/websocket 已经archived，不再维护！！！

---

**需要实现的功能：**

- [x] upgrade
- [x] read frame，加解码这个事情，不知道为什么老喜欢做。其实没啥好做的，就一个烦字，而且出问题很难定位。
- [x] write frame
- [ ] close

~~~shell
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
~~~

