package buffer

import (
	"bytes"
	"sync"
)

type Pool struct {
	sync.Pool
	SizePool   int
	SizeBuffer int
}

func NewBufferPool(size int, bufSize int) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(bytes.Buffer)
		}}
}
