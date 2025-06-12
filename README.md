im-fun
---


[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go](https://github.com/zhixunjie/im-fun/actions/workflows/go.yml/badge.svg)](https://github.com/zhixunjie/im-fun/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhixunjie/im-fun)](https://goreportcard.com/report/github.com/zhixunjie/im-fun)

IM PlayGround，Just For Fun 😄. Inspired by [Terry-Mao/goim](https://github.com/Terry-Mao/goim)

![flow.png](img/flow.png)


用途：学习 or 生产发布（给我个 star，让我知道项目对你有用吧！）

技术魔法：
1. 逻辑层设计：
   - 支持单聊/群聊；
   - 聊天记录的读写扩散解决方案；
   - 分库分表方式；
   - 消息推拉机制的实现；
   
4. 长链接设计：
   - 服务注册 / 服务发现
   - 单播/组播/广播的实现；
   - 缓冲池的内存分配、内存复用的实现；
   - 通信协议的设计、编码/解码的实现；
   - WebSocket 的协议解析；Upgrade 机制；
   - 心跳机制的设计；

