package scheduler

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func stubResolver(api hetzner.API) func(context.Context, int64) (hetzner.API, error) {
	return func(_ context.Context, _ int64) (hetzner.API, error) {
		return api, nil
	}
}

type recordingClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *recordingClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}
func (c *recordingClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func newStoreForScheduler(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestSchedulerTriggersRescaleOnWindowEntry(t *testing.T) {
	st := newStoreForScheduler(t)
	p, _ := st.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21", "cpx11"}, Mode: "scheduled", Timezone: "UTC",
	})
	_, _ = st.CreateWindow(srv.ID, store.Window{
		Label:      "all day",
		DaysOfWeek: 0b01111111, StartTime: "00:00", StopTime: "23:59", TargetType: "cpx21", Enabled: true,
	})

	api := hcloudmock.New()
	hserver := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(hserver)

	clk := &recordingClock{t: time.Date(2026, 6, 29, 0, 30, 0, 0, time.UTC)}

	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	done := make(chan struct{})
	go func() { sched.Run(); close(done) }()
	time.Sleep(200 * time.Millisecond)
	sched.Stop()
	<-done
}

func TestSchedulerAutoPromoteTriggersRescale(t *testing.T) {
	st := newStoreForScheduler(t)
	p, _ := st.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21", "cpx11"}, Mode: "auto_promote", Timezone: "UTC",
	})

	api := hcloudmock.New()
	hserver := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(hserver)

	clk := &recordingClock{t: time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)}

	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	// Set promote_state and tick manually
	ps := "promote_requested"
	srv.PromoteState = &ps
	if err := st.UpdateServer(*srv); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}

	// Manually call tick once
	sched.tick(srv.ID)

	// Allow a brief moment for events to flush
	time.Sleep(50 * time.Millisecond)

	events, err := st.ListEventsByServer(srv.ID, 10)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Kind != "rescale_up" || !events[0].OK {
		t.Fatalf("got event %+v", events[0])
	}
}

func seedSchedulerTestServer(t *testing.T, st *store.Store) *store.Server {
	t.Helper()
	p, err := st.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	srv, err := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "w", Label: "w",
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain: []string{"cpx31", "cpx11"},
		Mode:          "auto_promote", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	return srv
}

func TestScheduler_WritesSchedulerTickOnIdleAutoPromote(t *testing.T) {
	st := newStoreForScheduler(t)
	srv := seedSchedulerTestServer(t, st)

	api := hcloudmock.New()
	api.AddServer(&hetzner.Server{ID: srv.HCloudServerID, Name: srv.Name, ServerType: &hetzner.ServerType{Name: "cpx31"}})

	clk := &recordingClock{t: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}
	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	// Server at top + promote_requested → tick finds current == top → "ok_idle".
	ps := "promote_requested"
	srv.PromoteState = &ps
	if err := st.UpdateServer(*srv); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}
	sched.tick(srv.ID)

	events, err := st.ListEventsByServer(srv.ID, 10)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	var found *store.Event
	for _, e := range events {
		if e.Kind == "scheduler_tick" {
			found = e
			break
		}
	}
	if found == nil {
		t.Fatalf("expected scheduler_tick event, got %+v", events)
	}
	if found.Error != "ok_idle" {
		t.Fatalf("Error = %q, want ok_idle", found.Error)
	}
	if found.TriggeredBy != "scheduler" {
		t.Fatalf("TriggeredBy = %q, want scheduler", found.TriggeredBy)
	}
}

func TestScheduler_DebouncesTickHeartbeat(t *testing.T) {
	st := newStoreForScheduler(t)
	srv := seedSchedulerTestServer(t, st)

	api := hcloudmock.New()
	api.AddServer(&hetzner.Server{ID: srv.HCloudServerID, Name: srv.Name, ServerType: &hetzner.ServerType{Name: "cpx31"}})

	clk := &recordingClock{t: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}
	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	ps := "promote_requested"
	srv.PromoteState = &ps
	_ = st.UpdateServer(*srv)

	sched.tick(srv.ID)
	clk.advance(2 * time.Minute)
	sched.tick(srv.ID) // within debounce window — should NOT write again.
	clk.advance(4 * time.Minute)
	sched.tick(srv.ID) // past debounce window — should write.

	events, _ := st.ListEventsByServer(srv.ID, 10)
	count := 0
	for _, e := range events {
		if e.Kind == "scheduler_tick" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("scheduler_tick count = %d, want 2 (debounce should have suppressed one)", count)
	}
}

func TestScheduler_WritesTickOnLockContention(t *testing.T) {
	st := newStoreForScheduler(t)
	srv := seedSchedulerTestServer(t, st)

	api := hcloudmock.New()
	api.AddServer(&hetzner.Server{ID: srv.HCloudServerID, Name: srv.Name, ServerType: &hetzner.ServerType{Name: "cpx11"}})

	clk := &recordingClock{t: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}
	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	// Hold the action lock so AcquireAction returns false.
	acquired, err := st.AcquireAction(srv.ID, "rescale_to_cpx31", 30*time.Minute)
	if err != nil || !acquired {
		t.Fatalf("seed lock: acquired=%v err=%v", acquired, err)
	}

	ps := "promote_requested"
	srv.PromoteState = &ps
	_ = st.UpdateServer(*srv)
	sched.tick(srv.ID)

	events, _ := st.ListEventsByServer(srv.ID, 10)
	var found *store.Event
	for _, e := range events {
		if e.Kind == "scheduler_tick" {
			found = e
			break
		}
	}
	if found == nil || found.Error != "lock_contention" {
		t.Fatalf("expected scheduler_tick lock_contention, got %+v", events)
	}
}

func TestScheduler_WritesTickOnNoWindows(t *testing.T) {
	st := newStoreForScheduler(t)
	p, _ := st.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	srv, _ := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "w", Label: "w",
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain: []string{"cpx31"}, Mode: "scheduled", Timezone: "UTC",
	})

	api := hcloudmock.New()
	api.AddServer(&hetzner.Server{ID: srv.HCloudServerID, Name: srv.Name, ServerType: &hetzner.ServerType{Name: "cpx11"}})

	clk := &recordingClock{t: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}
	sched := New(st, stubResolver(api), clk, 50*time.Millisecond)
	sched.Add(srv.ID)

	sched.tick(srv.ID)

	events, _ := st.ListEventsByServer(srv.ID, 10)
	var found *store.Event
	for _, e := range events {
		if e.Kind == "scheduler_tick" {
			found = e
			break
		}
	}
	if found == nil || found.Error != "no_windows" {
		t.Fatalf("expected scheduler_tick no_windows, got %+v", events)
	}
}

func TestAttachLifecycle(t *testing.T) {
	st := newStoreForScheduler(t)
	p, _ := st.CreateProject("p", []byte("tok"), []byte("nonce12byts"))

	clk := &recordingClock{t: time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)}
	sched := New(st, stubResolver(hcloudmock.New()), clk, 50*time.Millisecond)

	hub := broadcast.NewHub[store.ServerLifecycleEvent]()
	st.SetServerLifecycleHub(hub)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() { _ = sched.Attach(ctx, hub); close(done) }()

	srv, _ := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "w", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC",
	})

	// Wait for the lifecycle hub to deliver "created" to Attach, which calls sched.Add.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		sched.mu.Lock()
		_, ok := sched.added[srv.ID]
		sched.mu.Unlock()
		if ok {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	sched.mu.Lock()
	_, ok := sched.added[srv.ID]
	sched.mu.Unlock()
	if !ok {
		t.Fatalf("scheduler did not Add the server after lifecycle broadcast")
	}

	cancel()
	sched.Stop()
	<-done
}
