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

# 优点说明

- Bufio ：
    - Bufio 复用了 Buffer，从而减少每个TCP的IO读写带来的Buffer GC。
- Buffer Pool：
    - 通过链表方式分配内存。
    - 每次内存增长时，会预先分配一大段内存再进行切分。
