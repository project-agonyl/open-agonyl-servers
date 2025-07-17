package shared

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// SafeQueue is a high-performance, goroutine-safe generic queue
type SafeQueue[T any] struct {
	// Ring buffer for lock-free operations
	buffer []atomic.Value
	head   uint64
	tail   uint64
	mask   uint64
	size   uint64

	// Synchronization
	notEmpty  chan struct{}
	notFull   chan struct{}
	closed    int32
	closeOnce sync.Once

	// Metrics for monitoring
	enqueueCount uint64
	dequeueCount uint64
	dropCount    uint64
}

// NewSafeQueue creates a new generic SafeQueue with the specified capacity
// Capacity must be a power of 2 for optimal performance
func NewSafeQueue[T any](capacity uint64) *SafeQueue[T] {
	// Ensure capacity is power of 2
	if capacity == 0 || (capacity&(capacity-1)) != 0 {
		panic("capacity must be a power of 2")
	}

	sq := &SafeQueue[T]{
		buffer:   make([]atomic.Value, capacity),
		mask:     capacity - 1,
		size:     capacity,
		notEmpty: make(chan struct{}, 1),
		notFull:  make(chan struct{}, 1),
	}

	// Initialize notFull channel
	select {
	case sq.notFull <- struct{}{}:
	default:
	}

	return sq
}

// Enqueue adds an item to the queue
// Returns true if successful, false if queue is full or closed
func (sq *SafeQueue[T]) Enqueue(item T) bool {
	if atomic.LoadInt32(&sq.closed) == 1 {
		return false
	}

	currentTail := atomic.LoadUint64(&sq.tail)
	nextTail := currentTail + 1

	// Check if queue is full
	if nextTail-atomic.LoadUint64(&sq.head) > sq.size {
		atomic.AddUint64(&sq.dropCount, 1)
		return false
	}

	// Try to claim the slot
	if !atomic.CompareAndSwapUint64(&sq.tail, currentTail, nextTail) {
		// Someone else got it, try again
		return sq.Enqueue(item)
	}

	// Store the item
	sq.buffer[currentTail&sq.mask].Store(item)
	atomic.AddUint64(&sq.enqueueCount, 1)

	// Signal that queue is not empty
	select {
	case sq.notEmpty <- struct{}{}:
	default:
	}

	return true
}

// TryEnqueue attempts to enqueue with a timeout
func (sq *SafeQueue[T]) TryEnqueue(item T, timeout time.Duration) bool {
	if sq.Enqueue(item) {
		return true
	}

	// If immediate enqueue failed, wait for space
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(time.Microsecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if sq.Enqueue(item) {
				return true
			}
		}
	}
}

// Dequeue removes and returns an item from the queue
// Returns the zero value and false if queue is empty
func (sq *SafeQueue[T]) Dequeue() (T, bool) {
	var zero T
	currentHead := atomic.LoadUint64(&sq.head)

	// Check if queue is empty
	if currentHead >= atomic.LoadUint64(&sq.tail) {
		return zero, false
	}

	// Try to claim the slot
	if !atomic.CompareAndSwapUint64(&sq.head, currentHead, currentHead+1) {
		// Someone else got it, try again
		return sq.Dequeue()
	}

	// Get the item
	slot := &sq.buffer[currentHead&sq.mask]
	item := slot.Load().(T)
	slot.Store(nil) // Clear reference for GC

	atomic.AddUint64(&sq.dequeueCount, 1)

	// Signal that queue is not full
	select {
	case sq.notFull <- struct{}{}:
	default:
	}

	return item, true
}

// DequeueBlocking blocks until an item is available or context is cancelled
func (sq *SafeQueue[T]) DequeueBlocking(ctx context.Context) (T, bool) {
	var zero T
	for {
		if item, ok := sq.Dequeue(); ok {
			return item, true
		}

		if atomic.LoadInt32(&sq.closed) == 1 {
			return zero, false
		}

		select {
		case <-ctx.Done():
			return zero, false
		case <-sq.notEmpty:
			// Continue to try dequeue
		}
	}
}

// DequeueTimeout attempts to dequeue with a timeout
func (sq *SafeQueue[T]) DequeueTimeout(timeout time.Duration) (T, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return sq.DequeueBlocking(ctx)
}

// Size returns the current number of items in the queue
func (sq *SafeQueue[T]) Size() uint64 {
	tail := atomic.LoadUint64(&sq.tail)
	head := atomic.LoadUint64(&sq.head)
	return tail - head
}

// IsEmpty returns true if the queue is empty
func (sq *SafeQueue[T]) IsEmpty() bool {
	return sq.Size() == 0
}

// IsFull returns true if the queue is full
func (sq *SafeQueue[T]) IsFull() bool {
	return sq.Size() >= sq.size
}

// Close closes the queue, preventing new enqueues
func (sq *SafeQueue[T]) Close() {
	sq.closeOnce.Do(func() {
		atomic.StoreInt32(&sq.closed, 1)
		close(sq.notEmpty)
		close(sq.notFull)
	})
}

// Stats returns queue statistics
func (sq *SafeQueue[T]) Stats() (enqueued, dequeued, dropped, current uint64) {
	return atomic.LoadUint64(&sq.enqueueCount),
		atomic.LoadUint64(&sq.dequeueCount),
		atomic.LoadUint64(&sq.dropCount),
		sq.Size()
}

// DrainAll removes all items from the queue and returns them
func (sq *SafeQueue[T]) DrainAll() []T {
	var items []T
	for {
		item, ok := sq.Dequeue()
		if !ok {
			break
		}
		items = append(items, item)
	}
	return items
}
