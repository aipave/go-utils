package gcache

import (
    "sync"
    "time"

    "github.com/RussellLuo/timingwheel"
)

// expireQueue Store "key expire automatically deleted after task"
type expireQueue[K comparable, V any] struct {
    tick      time.Duration
    wheelSize int64
    // Time round storage time due to delete task
    tw *timingwheel.TimingWheel

    mu     sync.Mutex
    timers map[K]*timingwheel.Timer
}

// newExpireQueue
///> Generate expiration queue object, the queue by time round, the elements in the queue according to the expiration time, regularly delete
func newExpireQueue[K comparable, V any](tick time.Duration, wheelSize int64) *expireQueue[K, V] {
    queue := &expireQueue[K, V]{
        tick:      tick,
        wheelSize: wheelSize,

        tw:     timingwheel.NewTimingWheel(tick, wheelSize),
        timers: make(map[K]*timingwheel.Timer),
    }

    // start a goroutine，handle expired entry
    go queue.tw.Start()
    return queue
}

// add Add time overdue tasks, each time the task executes, due to independent goroutine is carried out
func (q *expireQueue[K, V]) add(key K, expireTime time.Time, f func()) {
    q.mu.Lock()
    defer q.mu.Unlock()

    d := expireTime.Sub(currentTime())
    timer := q.tw.AfterFunc(d, q.task(key, f))
    q.timers[key] = timer

    return
}

// update key's tll
func (q *expireQueue[K, V]) update(key K, expireTime time.Time, f func()) {
    q.mu.Lock()
    defer q.mu.Unlock()

    if timer, ok := q.timers[key]; ok {
        timer.Stop()
    }

    d := expireTime.Sub(currentTime())
    timer := q.tw.AfterFunc(d, q.task(key, f))
    q.timers[key] = timer
}

// remove key
func (q *expireQueue[K, V]) remove(key K) {
    q.mu.Lock()
    defer q.mu.Unlock()

    if timer, ok := q.timers[key]; ok {
        timer.Stop()
        delete(q.timers, key)
    }
}

// clear
func (q *expireQueue[K, V]) clear() {
    q.tw.Stop()
    q.tw = timingwheel.NewTimingWheel(q.tick, q.wheelSize)
    q.timers = make(map[K]*timingwheel.Timer)

    // restart goroutine，handle expired entry
    go q.tw.Start()
}

// stop The pause time round run queue
func (q *expireQueue[K, V]) stop() {
    q.tw.Stop()
}

func (q *expireQueue[K, V]) task(key K, f func()) func() {
    return func() {
        f()
        q.mu.Lock()
        delete(q.timers, key)
        q.mu.Unlock()
    }
}
