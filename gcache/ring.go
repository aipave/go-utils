package gcache

import (
	"container/list"
	"sync"
)

// ringConsumer accept and consume data
type ringConsumer interface {
	push([]*list.Element) bool
}

// ringStrip ring buffer Metadata caching Get requests to batch update element is located in the LRU
type ringStripe struct {
	consumer ringConsumer
	data     []*list.Element
	capacity int
}

func newRingStripe(consumer ringConsumer, capacity int) *ringStripe {
	return &ringStripe{
		consumer: consumer,
		data:     make([]*list.Element, 0, capacity),
		capacity: capacity,
	}
}

// push Recorded by a visit to the element to the ring buffer
func (r *ringStripe) push(ele *list.Element) {
	r.data = append(r.data, ele)
	if len(r.data) >= r.capacity {
		if r.consumer.push(r.data) {
			r.data = make([]*list.Element, 0, r.capacity)
		} else {
			r.data = r.data[:0]
		}
	}
}

// Pooling ringStripes of ringBuffer allows multiple goroutines to write elements to
// the buffer without locking, which is more efficient than concurrent writes to the same channel.
// Additionally, objects in the pool will be automatically removed without any notification,
// which reduces cache's LRU operations and minimizes concurrent competition between Set/Expire/Delete operations.
// The cache doesn't need to be strictly LRU,
// actively discarding some access metadata to reduce concurrent competition and
// improve write efficiency is necessary
type ringBuffer struct {
	pool *sync.Pool
}

func newRingBuffer(consumer ringConsumer, capacity int) *ringBuffer {
	return &ringBuffer{
		pool: &sync.Pool{
			New: func() interface{} { return newRingStripe(consumer, capacity) },
		},
	}
}

// push Record a is access to the elements
func (b *ringBuffer) push(ele *list.Element) {
	ringStripe := b.pool.Get().(*ringStripe)
	ringStripe.push(ele)
	b.pool.Put(ringStripe)
}
