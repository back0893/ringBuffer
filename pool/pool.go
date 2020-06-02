package pool

import (
	"github.com/back0893/ringBuffer"
	"sync"
)

var DefaultPool = NewPool(1024)

func Get() *ringbuffer.RingBuffer {
	return DefaultPool.Get()
}
func Put(r *ringbuffer.RingBuffer) {
	DefaultPool.Put(r)
}

type RingBufferPool struct {
	pool *sync.Pool
}

func NewPool(initSize int) *RingBufferPool {
	return &RingBufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return ringbuffer.NewRingBuffer(initSize)
			},
		},
	}
}
func (p *RingBufferPool) Get() *ringbuffer.RingBuffer {
	return p.pool.Get().(*ringbuffer.RingBuffer)
}
func (p *RingBufferPool) Put(r *ringbuffer.RingBuffer) {
	p.pool.Put(r)
}
