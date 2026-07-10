package rescaler

import (
	"context"
	"errors"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

func TestIsTypeAvailable_AvailableInLocation(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "fsn1", true)
	srv := &hetzner.Server{
		Datacenter: &hetzner.Datacenter{
			Location: &hetzner.Location{Name: "fsn1"},
		},
	}
	ok, err := IsTypeAvailable(context.Background(), api, srv, "cpx11")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if !ok {
		t.Fatal("ok = false, want true")
	}
}

func TestIsTypeAvailable_UnavailableInLocation(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "fsn1", false)
	srv := &hetzner.Server{
		Datacenter: &hetzner.Datacenter{
			Location: &hetzner.Location{Name: "fsn1"},
		},
	}
	ok, err := IsTypeAvailable(context.Background(), api, srv, "cpx11")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if ok {
		t.Fatal("ok = true, want false")
	}
}

func TestIsTypeAvailable_NotInLocation(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "nbg1", true) // cpx11 only in nbg1, not fsn1
	srv := &hetzner.Server{
		Datacenter: &hetzner.Datacenter{
			Location: &hetzner.Location{Name: "fsn1"},
		},
	}
	ok, err := IsTypeAvailable(context.Background(), api, srv, "cpx11")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if ok {
		t.Fatal("ok = true, want false (type not in location)")
	}
}

func TestIsTypeAvailable_APIError(t *testing.T) {
	api := hcloudmock.New()
	api.SetGetServerTypeError("cpx11", errors.New("network blip"))
	srv := &hetzner.Server{
		Datacenter: &hetzner.Datacenter{
			Location: &hetzner.Location{Name: "fsn1"},
		},
	}
	_, err := IsTypeAvailable(context.Background(), api, srv, "cpx11")
	if err == nil {
		t.Fatal("err = nil, want non-nil (fail open is the caller's decision)")
	}
}

func TestIsTypeAvailable_ServerWithoutDatacenter(t *testing.T) {
	api := hcloudmock.New()
	srv := &hetzner.Server{} // no Datacenter
	ok, err := IsTypeAvailable(context.Background(), api, srv, "cpx11")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if ok {
		t.Fatal("ok = true, want false")
	}
	if api.GetServerTypeCallCount("cpx11") != 0 {
		t.Fatalf("GetServerType was called despite nil Datacenter; expected zero calls")
	}
}
