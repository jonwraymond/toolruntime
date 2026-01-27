package gvisor

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/toolruntime"
)

func TestBackendImplementsInterface(t *testing.T) {
	var _ toolruntime.Backend = (*Backend)(nil)
}

func TestBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendGVisor {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendGVisor)
	}
}

func TestBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.runscPath != "runsc" {
		t.Errorf("runscPath = %q, want %q", b.runscPath, "runsc")
	}
	if b.platform != "systrap" {
		t.Errorf("platform = %q, want %q", b.platform, "systrap")
	}
	if b.networkMode != "none" {
		t.Errorf("networkMode = %q, want %q", b.networkMode, "none")
	}
}

func TestBackendRequiresGateway(t *testing.T) {
	b := New(Config{})
	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "test",
		Gateway: nil,
	}
	_, err := b.Execute(ctx, req)
	if !errors.Is(err, toolruntime.ErrMissingGateway) {
		t.Errorf("Execute() without gateway error = %v, want %v", err, toolruntime.ErrMissingGateway)
	}
}
