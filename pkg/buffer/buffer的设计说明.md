# 1. 数据结构

1、Pool & Buffer：缓冲池子

~~~
Pool -> Buffer
     		-> Get -> get a Buffer from the Pool
     		-> Put -> put back a Buffer to the Pool
~~~

2、Bufio：[改造版的Bufio](./bufio/bufio(缓冲区读写-增强).md)。

3、Pool Hash：使用多个池子，基于取余进行池子分配。减少池子的互斥情况。

---

**相互调用关系：**

1. Pool 管理 Buffer（池子里面的一个单元，就是一个Buffer对象）。
2. Bufio相当于让TCP Reader(conn)的读写带上了缓冲区（相当于C语言的标准IO函数，一大块用户缓冲区），从而减少conn的系统调用Read/Write的次数。
  - 由于Bufio的用户缓冲区会每个TCP连接都带上，如果频繁进行创建和销毁，申请内存和GC都要消耗性能的。
  - 所以，用户缓冲区的内存交由Buffer Pool去管理。

# 2. 优点说明

- Bufio ：复用了 Buffer Pool里面的Buffer，从而减少每个TCP的IO读写带来的Buffer GC。
- Buffer Pool：
    - 通过链表方式分配内存。
    - 当发现Buffer Pool没有Buffer时，需要预先分配一大段内存再进行切分（批量创建buffer）。
      - 相对于golang自带的sync.Pool， 好处就是批量New，而不是一个个去New。
