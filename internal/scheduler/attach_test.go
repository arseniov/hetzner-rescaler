package scheduler

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestAttach_RuntimeCreatedServerGetsTick(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	p, err := st.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	lifecycleHub := broadcast.NewHub[store.ServerLifecycleEvent]()
	st.SetServerLifecycleHub(lifecycleHub)

	// recordingClock at 12:00 UTC; window is 01:00-02:00 (clearly outside).
	// This keeps the test fully deterministic without needing a mock hcloud.
	clk := &recordingClock{t: time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)}
	sched := New(st,
		func(_ context.Context, _ int64) (hetzner.API, error) { return hcloudmock.New(), nil },
		clk, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() { _ = sched.Attach(ctx, lifecycleHub); close(done) }()

	srv, err := st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "w", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "scheduled", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	_, _ = st.CreateWindow(srv.ID, store.Window{
		Label:      "small hours",
		DaysOfWeek: 0b01111111, StartTime: "01:00", StopTime: "02:00", TargetType: "cpx21", Enabled: true,
	})

	// Give Attach's subscriber goroutine a moment to consume the
	// "created" event and call sched.Add(srv.ID).
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		sched.mu.Lock()
		_, added := sched.added[srv.ID]
		sched.mu.Unlock()
		if added {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Run one tick: out-of-window → writeTickSummary("ok_idle") emits a
	// scheduler_tick event, proving the scheduler is alive and reachable
	// from the runtime-created server.
	sched.tick(srv.ID)

	events, err := st.ListEventsByServer(srv.ID, 100)
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
		t.Fatalf("Error = %q, want %q", found.Error, "ok_idle")
	}

	cancel()
	sched.Stop()
	<-done
}