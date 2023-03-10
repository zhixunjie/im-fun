package buffer

type Options struct {
	ReadPoolOption  PoolOptions `yaml:"readPoolOption"`
	WritePoolOption PoolOptions `yaml:"writePoolOption"`
}

type PoolOptions struct {
	PoolNum  int `yaml:"poolNum"`  // 池子的个数
	BatchNum int `yaml:"batchNum"` // 池子创建Buffer时批量创建的个数
	BufSize  int `yaml:"bufSize"`  // 每个Buffer的字节数
}

type PoolHash struct {
	Readers []Pool
	Writers []Pool
	options *Options
}

func NewPoolHash(config *Options) PoolHash {
	var hash PoolHash
	hash.options = config

	// new reader pool
	var option PoolOptions
	option = hash.options.ReadPoolOption
	hash.Readers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.Readers[i].Init(option.BufSize, option.BatchNum)
	}

	// new writer pool
	option = hash.options.WritePoolOption
	hash.Writers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.Writers[i].Init(option.BufSize, option.BatchNum)
	}

	return hash
}

// ReaderPool get a reader memory buffer.
func (pool *PoolHash) ReaderPool(rn int) *Pool {
	return &(pool.Readers[rn%pool.options.ReadPoolOption.PoolNum])
}

// WriterPool get a writer memory buffer pool.
func (pool *PoolHash) WriterPool(rn int) *Pool {
	return &(pool.Writers[rn%pool.options.WritePoolOption.PoolNum])
}
