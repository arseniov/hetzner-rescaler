package hetzner

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ServerType is a small re-export of hcloud.ServerType to keep the engine
// decoupled from the SDK's exact struct shape (mostly — we still use it
// directly for fields like .Name).
type ServerType = hcloud.ServerType

// Server is a small re-export of hcloud.Server.
type Server = hcloud.Server

// Action is a small re-export of hcloud.Action.
type Action = hcloud.Action

// Datacenter is a small re-export of hcloud.Datacenter. (Deprecated in the
// Hetzner Cloud API, but still populated on hcloud.Server for now.)
type Datacenter = hcloud.Datacenter

// Location is a small re-export of hcloud.Location.
type Location = hcloud.Location

// Action status constants, re-exported from hcloud.ActionStatus* for use in
// tests and in the rescaler code that polls action status.
const (
	ActionStatusRunning = hcloud.ActionStatusRunning
	ActionStatusSuccess = hcloud.ActionStatusSuccess
	ActionStatusError   = hcloud.ActionStatusError
)
