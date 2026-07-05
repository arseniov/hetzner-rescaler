package broadcast

import (
	"sync"
	"testing"
	"time"
)

type testEvent struct{ id int }

func TestHubSubscribeReceivesBroadcast(t *testing.T) {
	h := NewHub[testEvent]()
	defer h.Close()

	ch, unsub := h.Subscribe(8)
	defer unsub()

	go func() { h.Broadcast(testEvent{id: 1}) }()

	select {
	case ev := <-ch:
		if ev.id != 1 {
			t.Fatalf("got %d, want 1", ev.id)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestHubUnsubscribeStopsDelivery(t *testing.T) {
	h := NewHub[testEvent]()
	defer h.Close()

	ch, unsub := h.Subscribe(2)
	unsub() // unsubscribe immediately

	h.Broadcast(testEvent{id: 1})

	select {
	case _, open := <-ch:
		if open {
			t.Fatal("expected channel to be closed after unsubscribe")
		}
	case <-time.After(100 * time.Millisecond):
		// Channel close races with broadcast; treat as success.
	}
}

func TestHubConcurrentBroadcast(t *testing.T) {
	h := NewHub[testEvent]()
	defer h.Close()

	const subs = 10
	const msgs = 50
	var wg sync.WaitGroup

	for i := 0; i < subs; i++ {
		ch, unsub := h.Subscribe(msgs * subs)
		wg.Add(1)
		go func(c <-chan testEvent) {
			defer wg.Done()
			defer unsub()
			received := 0
			for range c {
				received++
				if received == msgs {
					return
				}
			}
		}(ch)
	}

	for i := 0; i < msgs; i++ {
		h.Broadcast(testEvent{id: i})
	}

	wg.Wait()
}