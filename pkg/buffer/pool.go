package buffer

import "sync"

// Pool：自己编写的BufferPool
// TODO try sth lock free，like: sync.Pool

type Pool struct {
	lock     sync.Mutex
	free     *Buffer // check this detail in BatchNew
	bufSize  int     // each Buffer Size
	batchNum int     // create Buffer continuously
}

// Init init Pool
func (pool *Pool) Init(bufSize, batchNum int) {
	pool.bufSize = bufSize
	pool.batchNum = batchNum
	pool.BatchNew()
}

// BatchNew Buffer
// free LinkList:  free buffer 1 ->  free buffer 2 -> free buffer 3
func (pool *Pool) BatchNew() {
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
		pool.BatchNew()
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
