// Package broadcast provides a tiny in-process pub/sub hub.
//
// One Hub per process is typically enough; it's safe for concurrent use.
// Subscribers receive every value sent to Broadcast while they are
// subscribed; if a subscriber's buffered channel is full, the broadcast
// is dropped for that subscriber (the next broadcast may succeed).
package broadcast

import "sync"

type Hub[T any] struct {
	mu   sync.Mutex
	subs map[*sub[T]]struct{}
}

type sub[T any] struct {
	ch chan T
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{subs: map[*sub[T]]struct{}{}}
}

// Subscribe registers a new subscriber and returns its channel and an
// unsubscribe function. Buffer size hints at how many events may queue
// before delivery is dropped for this subscriber.
func (h *Hub[T]) Subscribe(buffer int) (<-chan T, func()) {
	s := &sub[T]{ch: make(chan T, buffer)}
	h.mu.Lock()
	h.subs[s] = struct{}{}
	h.mu.Unlock()

	var once sync.Once
	unsub := func() {
		once.Do(func() {
			h.mu.Lock()
			delete(h.subs, s)
			h.mu.Unlock()
			close(s.ch)
		})
	}
	return s.ch, unsub
}

// Broadcast delivers ev to every current subscriber. Subscribers with
// full buffers miss this event; we don't block on slow subscribers.
func (h *Hub[T]) Broadcast(ev T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for s := range h.subs {
		select {
		case s.ch <- ev:
		default:
		}
	}
}

// Close releases hub resources. Subscribers must consume until their
// channels close.
func (h *Hub[T]) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for s := range h.subs {
		close(s.ch)
		delete(h.subs, s)
	}
}