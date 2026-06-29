package scheduler

import (
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestIsInWindow_ExactMatch(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: true}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, loc)

	got, target, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if !got {
		t.Fatal("expected in-window at 12:00")
	}
	if target != "" {
		// TargetType is empty in this window struct; that's fine.
	}
}

func TestIsInWindow_BeforeStart(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: true}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 8, 59, 0, 0, loc)

	got, _, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if got {
		t.Fatal("expected not-in-window at 08:59")
	}
}

func TestIsInWindow_AfterStop(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: true}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 19, 1, 0, 0, loc)

	got, _, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if got {
		t.Fatal("expected not-in-window at 19:01")
	}
}

func TestIsInWindow_WrongDay(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: true} // Mon-Fri
	loc, _ := time.LoadLocation("UTC")
	// 2026-06-28 is a Sunday
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, loc)

	got, _, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if got {
		t.Fatal("expected not-in-window on Sunday for Mon-Fri window")
	}
}

func TestIsInWindow_TimezoneAware(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b0111111, Enabled: true}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, loc)

	got, _, err := EvaluateWindows([]store.Window{w}, now, "Europe/Rome")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if !got {
		t.Fatal("expected in-window in Europe/Rome at 14:00 local")
	}
}

func TestIsInWindow_DisabledWindowIgnored(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: false}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, loc)

	got, _, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if got {
		t.Fatal("disabled window should not match")
	}
}

func TestEvaluateWindowsReturnsTargetType(t *testing.T) {
	w := store.Window{StartTime: "09:00", StopTime: "19:00", DaysOfWeek: 0b00111110, Enabled: true, TargetType: "cpx31"}
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, loc)

	got, target, err := EvaluateWindows([]store.Window{w}, now, "UTC")
	if err != nil {
		t.Fatalf("EvaluateWindows: %v", err)
	}
	if !got || target != "cpx31" {
		t.Fatalf("got in-window=%v target=%q, want true cpx31", got, target)
	}
}