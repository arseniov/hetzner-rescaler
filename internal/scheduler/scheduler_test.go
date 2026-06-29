package scheduler

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

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

	sched := New(st, api, clk, 50*time.Millisecond)
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

	sched := New(st, api, clk, 50*time.Millisecond)
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