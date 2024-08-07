# 参考

- https://github.com/halfrost/Halfrost-Field/blob/master/contents/Protocol/WebSocket.md
- https://en.wikipedia.org/wiki/WebSocket

# overview

The handshake from the client looks as follows:

~~~shell
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

~~~shell
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat
~~~

# ------ 客户端 ------ 

- PHP：https://www.kancloud.cn/zhixunjie/swoole_websocket/363058
- GO：https://github.com/zhixunjie/im-fun/tree/master/benchmarks/client/websocket

# ------ 服务端 ------ 

# 1. 核心功能清单

> 自己实现一些功能试试：https://github.com/zhixunjie/im-fun/tree/master/pkg/websocket

RFC：

- https://datatracker.ietf.org/doc/html/rfc6455
- https://datatracker.ietf.org/doc/html/rfc6455#section-1.3 upgrade
- https://datatracker.ietf.org/doc/html/rfc6455#section-5.2 frame

---

核心功能清单：

- [x] upgrade

- [x] frame：read frame、write frame

  > 加解码的事情，就一个烦字，出问题很难定位。

- [x] close

- [ ] heartbeat

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

# 2. 调研：第三方库

1. [gorilla](research/websocket库_gorilla.md)  

   > 历史悠久，star 多（推荐使用🔥🔥🔥）

2. [nhooyr](research/websocket库_nhooyr.md) ：比较新的包，star 少。

---

其他市面上 star 较多的库：

1. https://github.com/lxzan/gws

   > GWS (Go WebSocket) is a very simple, fast, reliable and feature-rich WebSocket implementation written in Go.

## 2.1 gorilla vs nhooyr

> 无论是[gorilla](https://github.com/gorilla/websocket/blob/master/examples/echo/server.go)，还是[nhooyr](https://github.com/nhooyr/websocket)，使用时都是基于net/http做上层的封装。

- gorilla：
  - 2023-08：最新消息，gorilla又回来了，[重新开始维护](https://github.com/gorilla#gorilla-toolkit)。
  - 2023-05：目前已经不再维护（archived），[作者给出的原因](https://github.com/gorilla#gorilla-toolkit)是该类库已经足够稳定，并且没有更多改进的空间了。
- nhooyr：比较容易看懂，但是功能还不够全面。

**最终结论**：

- 如果追求更全面的功能，可以看看gorilla（由于它又恢复维护了，所以重新推荐使用🔥🔥🔥）。
- 使用nhooyr，比较容易看懂，但是功能还不够全面。
