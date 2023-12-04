> IM PlayGround，Just For Fun 😄

数据流：

![image-20230305161416634](https://typroa-jasonzhi.oss-cn-guangzhou.aliyuncs.com/imgs/image-20230305161416634.png)

组件图：

![image-20230624000404596](https://typroa-jasonzhi.oss-cn-guangzhou.aliyuncs.com/imgs/image-20230624000404596.png)

# 技术魔法

1. 消息系统本身的设计：

    - timeline设计、读写扩散、message的session表的设计（各类id、分表分库）。
    - 功能点实现：消息发送、消息拉新、消息顺序、多端同步、消息未读数。
    - 长链接消息发送：单播（用户）、组播（房间）、广播（全部用户）。

- 消息及时性保证、消息的顺序保证。
3. [Bufer 设计](pkg/buffer/buffer的设计说明.md)： 
    - Buffer Pool：复用Buffer，内存复用，减少GC。
    - Bufio魔改：复用Buffer（Buffer来自Buffer Pool）。
      - 支持Peek：写入时复用写Buffer
      - 支持Pop：读取时复用Buffer。
4. Ring：环形数组

    - 复用proto，内存复用，减少GC。

    - 同时，用于限流，限制读写的频率。
5. Bucket设计、Round设计（Hash Buffer Pool、Hash Timer）
    - 协程分流：一个 TCP 连接对应一个协程。
        - Bucket、Buffer Pool、Timer的操作中存在很多加锁方法（锁冲突）。
        - 为了减少锁冲突，通过Hash的方式，不同协程协程会分配到不同的Bucket、Buffer Pool、Timer。
6. 哈希表、链表、数组的选择取舍。
    - Room里面的Channel，通过链表的形式存放，目的是增加/删除操作方便。
        - Room：ch1 -> ch2 -> ch3
    - Bucket里面的Channel，通过哈希表的形式存放，目的是快速判断某个Channel是否存在。
        - Bucket：UserKey => ch
    - 数组与链表：查找和增删操作的取舍，哈希表和链表：O(1)查找和顺序遍历的取舍。
        - 顺序遍历：数组、链表。
        - 快速查找：数组、哈希表。
        - 快速增删：哈希表、链表。如果增删最后的元素，数组也很快。
7. 中间层Job：
    - 使用Kafka，提高消息发送效率，异步发送消息。
    - 单独的Job进行Kafka消息，迭代升级时不会影响到comet提供的长连接服务。
8. 网络模型：协程池Accept（每个协程里面Accept）、后期可以改造为Reactor模型。

9. [WebSocket](pkg/websocket/docs/websocket-技术文档.md)： 
    - Upgrade机制：握手机制，发送什么样的数据包。
    - 服务端读写数据时编码和解码。
10. 服务注册、服务发现。
11. 分布式部署、大型IM系统。
13. 性能测试、性能调优、内存泄露。

# TCP的数据包格式

单个Proto的格式：由两部分组成：Header + Body。

- Package Length：整个包的长度（Header + Body）。
- Header Length：Header部分的长度。
- Protocol Version：数据包格式的版本。
- Operation：操作类型（对应不同类型操作的消息），如：接收批量消息、心跳请求/响应等。
- Sequence Id：序列号
- Body：请求体，根据业务需求，可以自定义不同的请求体内容。

![image-20230418112113592](https://typroa-jasonzhi.oss-cn-guangzhou.aliyuncs.com/imgs/image-20230418112113592.png)

批量Proto的格式：

![image-20230623235943605](https://typroa-jasonzhi.oss-cn-guangzhou.aliyuncs.com/imgs/image-20230623235943605.png)