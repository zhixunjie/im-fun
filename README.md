> IM PlayGround，Just For Fun 😄

![image-20230305161416634](.README.assets/image-20230305161416634.png)

技术魔法：
1. 消息系统本身的设计：
    - timeline设计、读写扩散、message的session表的设计（各类id、分表分库）。
    - 功能点实现：消息发送、消息拉新、消息顺序、多端同步、消息未读数。
    - 消息及时性保证、消息的顺序保证。
2. Buffer Pool：复用Buffer，内存复用，减少GC。
3. Bufio魔改：复用Buffer（Buffer来自Buffer Pool）。
    - 支持Peek：写入时复用写Buffer
    - 支持Pop：读取时复用Buffer。
4. Ring：环形数组，复用proto，内存复用，减少GC。同时，用于限流，限制读写的频率。
5. Bucket设计、Round设计（Hash Buffer Pool、Hash Timer）
    - 协程分流：
        - Bucket、Buffer Pool、Timer的操作中存在很多加锁方法。
        - 通过Hash的方式，不同协程协程分配到不同的Bucket，从而减少冲突。
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
8. 网络模型：协程池Accept（每个协程里面Accept）、Reactor模型。
9. WebSocket：
    - Upgrade机制：握手机制，发送什么样的数据包。
    - 服务端读写数据时编码和解码。
10. 长链接消息发送：单播（用户）、组播（房间）、广播（全部用户）。
11. 服务注册、服务发现。
12. 分布式部署、大型IM系统。
13. 性能测试、性能调优、内存泄露。