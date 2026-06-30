package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// itoa is a tiny helper used by tests to build path segments.
func itoa(i int64) string { return strconv.FormatInt(i, 10) }

// recorder runs an HTTP request through the handler and returns the
// response recorder. Lives here so every test file in this package can
// share it without redefining.
func recorder(t *testing.T, h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

// fakeHetzner is a tiny in-memory stub satisfying hetzner.API for tests.
// It returns pre-programmed responses and nil for the SDK-pointer methods
// that handlers in later tasks will exercise.
type fakeHetzner struct {
	servers []*hetzner.Server
	types   []*hetzner.ServerType
	actions []*hetzner.Action
}

func (f *fakeHetzner) ListServers(ctx context.Context) ([]*hetzner.Server, error) {
	return f.servers, nil
}

func (f *fakeHetzner) GetServer(ctx context.Context, id int) (*hetzner.Server, error) {
	for i := range f.servers {
		if f.servers[i].ID == id {
			return f.servers[i], nil
		}
	}
	return nil, nil
}

func (f *fakeHetzner) ShutdownServer(ctx context.Context, srv *hetzner.Server) (*hetzner.Action, error) {
	return nil, nil
}

func (f *fakeHetzner) ChangeServerType(ctx context.Context, srv *hetzner.Server, t *hetzner.ServerType) (*hetzner.Action, error) {
	return nil, nil
}

func (f *fakeHetzner) PowerOnServer(ctx context.Context, srv *hetzner.Server) (*hetzner.Action, error) {
	return nil, nil
}

func (f *fakeHetzner) ListServerTypes(ctx context.Context) ([]*hetzner.ServerType, error) {
	return f.types, nil
}

func (f *fakeHetzner) GetServerType(ctx context.Context, name string) (*hetzner.ServerType, error) {
	for i := range f.types {
		if f.types[i].Name == name {
			return f.types[i], nil
		}
	}
	return nil, nil
}

func (f *fakeHetzner) GetAction(ctx context.Context, id int) (*hetzner.Action, error) {
	return nil, nil
}
