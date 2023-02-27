package buffer

type Options struct {
	r PoolOptions
	w PoolOptions
}

type PoolOptions struct {
	PoolNum  int
	BatchNum int
	BufSize  int
}

type PoolHash struct {
	readers []Pool
	writers []Pool
	options Options
}

func NewPoolHash(config Options) *PoolHash {
	hash := new(PoolHash)
	hash.options = config

	// new reader pool
	var option PoolOptions
	option = hash.options.r
	hash.readers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.readers[i].Init(option.BufSize, option.BatchNum)
	}

	// new writer pool
	option = hash.options.w
	hash.writers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.writers[i].Init(option.BufSize, option.BatchNum)
	}

	return hash
}

// Reader get a reader memory buffer.
func (pool *PoolHash) Reader(rn int) *Pool {
	return &(pool.readers[rn%pool.options.r.PoolNum])
}

// Writer get a writer memory buffer pool.
func (pool *PoolHash) Writer(rn int) *Pool {
	return &(pool.writers[rn%pool.options.w.PoolNum])
}
