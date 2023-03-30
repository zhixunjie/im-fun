# 数据结构

Pool & Buffer：缓冲池子

~~~
Pool -> Buffer
     -> Get -> get a Buffer from the Pool
     -> Put -> put back a Buffer to the Pool
~~~

Bufio：改造版的Bufio

- first, set a Reader.
    - note：Reader is the TCP/WebSocket connection.
- second, set Reader's Buffer.
    - get Buffer from the Pool for reusing the Buffer. It's good for reducing GC.

Pool Hash：

- 使用多个池子，基于取余进行池子分配。减少池子的互斥情况。

---

**相互调用关系：**

- Pool 管理 Buffer。
- Bufio的缓冲区，相当于用户缓冲区。TCP Reader(conn)的缓冲区，相当于内核缓冲区。
- Bufio相当于让TCP Reader(conn)的读写带上了缓冲区，从而减少conn的Read/Write调用次数。
  - 由于Bufio的缓冲区会每个TCP连接都带上，如果频繁进行创建和销毁，申请内存和GC都要消耗性能的。
  - 所以，缓冲区的内存交由Buffer Pool去管理。

# 优点说明

- Bufio ：
    - Bufio 复用了 Buffer，从而减少每个TCP的IO读写带来的Buffer GC。
- Buffer Pool：
    - 通过链表方式分配内存。
    - 当发现内存池子没有缓冲时，需要预先分配一大段内存再进行切分（批量创建buffer）。
      - 相对于golang自带的sync.Pool， 好处就是批量New，而不是一个个去New。
