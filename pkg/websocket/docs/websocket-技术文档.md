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

RFC：

- https://datatracker.ietf.org/doc/html/rfc6455
- https://datatracker.ietf.org/doc/html/rfc6455#section-5.2 frame
- https://datatracker.ietf.org/doc/html/rfc6455#section-1.3 upgrade

**我的实现：**

- https://github.com/zhixunjie/im-fun/tree/master/pkg/websocket

## 2.1 功能清单

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

## 2.2 第三方库

> 用于参考 or 直接拿来用

- [nhooyr](research/websocket库_nhooyr.md)
- [gorilla](research/websocket库_gorilla.md)

对比说明：

- 推荐使用nhooyr，比较容易看懂！！！
- gorilla属于比较旧的包，目前已经不再维护，作者给出的原因是该类库已经足够稳定，并且没有更多改进的空间了。
