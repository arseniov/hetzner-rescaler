package hetzner

import "errors"

// ErrUnavailable is returned (wrapped) when a target server type is out of
// stock. The scheduler and rescaler match on errors.Is(err, ErrUnavailable).
var ErrUnavailable = errors.New("hetzner: target server type unavailable")

// IsUnavailable reports whether err is or wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }
