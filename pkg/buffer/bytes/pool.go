package bytes

import (
	"sync"
)

// Pool：自己编写的BufferPool
// TODO try sth lock free，like: sync.Pool

type Pool struct {
	lock     sync.Mutex
	free     *Buffer // check this detail in batchNew
	bufSize  int     // 每个 Buffer 的大小
	batchNum int     // 连续创建一批数量的 Buffer
}

// Init init Pool
func (pool *Pool) Init(bufSize, batchNum int) {
	pool.bufSize = bufSize
	pool.batchNum = batchNum
	pool.batchNew()
}

// 批量创建一批 Buffer，并将其链接起来
// - free LinkList:  free buffer 1 ->  free buffer 2 -> free buffer 3
func (pool *Pool) batchNew() {
	bufSize := pool.bufSize
	batchNum := pool.batchNum
	bfArr := make([]Buffer, batchNum)
	byArr := make([]byte, batchNum*bufSize)

	// begin to traverse
	pool.free = &bfArr[0]
	p := &bfArr[0]
	for i := 1; i < batchNum; i++ {
		p.buf = byArr[(i-1)*bufSize : i*bufSize]
		p.next = &bfArr[i]
		p = p.next
	}
	p.buf = byArr[(batchNum-1)*bufSize : batchNum*bufSize]
	p.next = nil
}

// Get a Buffer from Pool（free LinkList）
func (pool *Pool) Get() (b *Buffer) {
	pool.lock.Lock()
	if b = pool.free; b == nil {
		pool.batchNew()
		b = pool.free
	}
	pool.free = b.next
	pool.lock.Unlock()

	return
}

// Put back a Buffer to the Pool
func (pool *Pool) Put(b *Buffer) {
	pool.lock.Lock()
	b.next = pool.free
	pool.free = b
	pool.lock.Unlock()
}
