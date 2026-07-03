package scheduler

import (
	"fmt"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// EvaluateWindows returns (inWindow, targetType, err). targetType is the
// TargetType of the first matching enabled window, or "" if none match.
//
// DaysOfWeek is a bitmask indexed by Go's time.Weekday: bit 0 = Sunday,
// bit 1 = Monday, ..., bit 6 = Saturday. So a "Mon-Fri" mask is 0b00111110.
//
// A window is "in" when the local time in `timezone` is on a day whose bit
// is set, AND localTime is in [StartTime, StopTime). The StopTime is treated
// as exclusive.
func EvaluateWindows(windows []store.Window, now time.Time, timezone string) (bool, string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return false, "", fmt.Errorf("scheduler: load timezone %q: %w", timezone, err)
	}

	local := now.In(loc)

	for _, w := range windows {
		if !w.Enabled {
			continue
		}

		dayBit := 1 << uint(local.Weekday())
		if w.DaysOfWeek&dayBit == 0 {
			continue
		}

		start, err := parseHHMM(w.StartTime)
		if err != nil {
			return false, "", fmt.Errorf("scheduler: parse StartTime %q: %w", w.StartTime, err)
		}
		stop, err := parseHHMM(w.StopTime)
		if err != nil {
			return false, "", fmt.Errorf("scheduler: parse StopTime %q: %w", w.StopTime, err)
		}

		minutes := local.Hour()*60 + local.Minute()
		if minutes < start || minutes >= stop {
			continue
		}

		return true, w.TargetType, nil
	}

	return false, "", nil
}

// parseHHMM parses an "HH:MM" string into minutes since midnight.
func parseHHMM(s string) (int, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}