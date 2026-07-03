// Package hetzner wraps the Hetzner Cloud SDK behind a small interface.
// The engine depends only on this interface; the SDK is an implementation
// detail. Tests use internal/hcloudmock.
package hetzner

import (
	"context"
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// API is the small set of Hetzner operations the engine needs.
type API interface {
	// Server
	GetServer(ctx context.Context, id int) (*Server, error)
	ListServers(ctx context.Context) ([]*Server, error)
	ShutdownServer(ctx context.Context, srv *Server) (*Action, error)
	ChangeServerType(ctx context.Context, srv *Server, target *ServerType) (*Action, error)
	PowerOnServer(ctx context.Context, srv *Server) (*Action, error)

	// ServerType
	ListServerTypes(ctx context.Context) ([]*ServerType, error)
	GetServerType(ctx context.Context, name string) (*ServerType, error)

	// Action
	GetAction(ctx context.Context, id int) (*Action, error)
}

// realAPI wraps *hcloud.Client.
type realAPI struct {
	c *hcloud.Client
}

// NewClient returns a real API implementation for the given token.
func NewClient(token string) (API, error) {
	if token == "" {
		return nil, errors.New("hetzner: empty token")
	}
	return &realAPI{c: hcloud.NewClient(hcloud.WithToken(token))}, nil
}

func (a *realAPI) GetServer(ctx context.Context, id int) (*Server, error) {
	s, _, err := a.c.Server.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("hetzner: get server: %w", err)
	}
	if s == nil {
		return nil, fmt.Errorf("hetzner: server %d not found", id)
	}
	return s, nil
}

func (a *realAPI) ListServers(ctx context.Context) ([]*Server, error) {
	servers, err := a.c.Server.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("hetzner: list servers: %w", err)
	}
	return servers, nil
}

func (a *realAPI) ShutdownServer(ctx context.Context, srv *Server) (*Action, error) {
	act, _, err := a.c.Server.Shutdown(ctx, srv)
	if err != nil {
		return nil, fmt.Errorf("hetzner: shutdown: %w", err)
	}
	return act, nil
}

func (a *realAPI) ChangeServerType(ctx context.Context, srv *Server, target *ServerType) (*Action, error) {
	act, _, err := a.c.Server.ChangeType(ctx, srv, hcloud.ServerChangeTypeOpts{
		ServerType:  target,
		UpgradeDisk: false,
	})
	if err != nil {
		return nil, fmt.Errorf("hetzner: change type: %w", err)
	}
	return act, nil
}

func (a *realAPI) PowerOnServer(ctx context.Context, srv *Server) (*Action, error) {
	act, _, err := a.c.Server.Poweron(ctx, srv)
	if err != nil {
		return nil, fmt.Errorf("hetzner: poweron: %w", err)
	}
	return act, nil
}

func (a *realAPI) ListServerTypes(ctx context.Context) ([]*ServerType, error) {
	types, err := a.c.ServerType.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("hetzner: list server types: %w", err)
	}
	return types, nil
}

func (a *realAPI) GetServerType(ctx context.Context, name string) (*ServerType, error) {
	t, _, err := a.c.ServerType.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("hetzner: get server type: %w", err)
	}
	if t == nil {
		return nil, fmt.Errorf("hetzner: server type %q not found", name)
	}
	return t, nil
}

func (a *realAPI) GetAction(ctx context.Context, id int) (*Action, error) {
	act, _, err := a.c.Action.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("hetzner: get action: %w", err)
	}
	if act == nil {
		return nil, fmt.Errorf("hetzner: action %d not found", id)
	}
	return act, nil
}
