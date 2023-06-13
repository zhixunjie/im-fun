package bytes

// 利用Hash算法，均摊池子的请求流量

type Hash struct {
	Readers []Pool
	Writers []Pool
	options *Options
}

func NewHash(config *Options) Hash {
	var hash Hash
	hash.options = config

	// new reader pool
	var option PoolOptions
	option = hash.options.ReadPool
	hash.Readers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.Readers[i].Init(option.BufSize, option.BatchNum)
	}

	// new writer pool
	option = hash.options.WritePool
	hash.Writers = make([]Pool, option.PoolNum)
	for i := 0; i < option.PoolNum; i++ {
		hash.Writers[i].Init(option.BufSize, option.BatchNum)
	}

	return hash
}

// ReaderPool get a reader memory buffer.
func (pool *Hash) ReaderPool(rn int) *Pool {
	return &(pool.Readers[rn%pool.options.ReadPool.PoolNum])
}

// WriterPool get a writer memory buffer pool.
func (pool *Hash) WriterPool(rn int) *Pool {
	return &(pool.Writers[rn%pool.options.WritePool.PoolNum])
}

type Options struct {
	ReadPool  PoolOptions `yaml:"readPoolOption"`
	WritePool PoolOptions `yaml:"writePoolOption"`
}

type PoolOptions struct {
	PoolNum  int `yaml:"poolNum"`  // 池子的个数
	BatchNum int `yaml:"batchNum"` // 池子创建Buffer时批量创建的个数
	BufSize  int `yaml:"bufSize"`  // 每个Buffer的字节数
}
