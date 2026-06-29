package hetzner

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// ServerType is a small re-export of hcloud.ServerType to keep the engine
// decoupled from the SDK's exact struct shape (mostly — we still use it
// directly for fields like .Name).
type ServerType = hcloud.ServerType

// Server is a small re-export of hcloud.Server.
type Server = hcloud.Server

// Action is a small re-export of hcloud.Action.
type Action = hcloud.Action
