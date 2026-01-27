package containerd

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/toolruntime"
)

func TestBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*Backend)(nil)
}

func TestBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendContainerd {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendContainerd)
	}
}

func TestBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.imageRef != "toolruntime-sandbox:latest" {
		t.Errorf("imageRef = %q, want %q", b.imageRef, "toolruntime-sandbox:latest")
	}
	if b.namespace != "default" {
		t.Errorf("namespace = %q, want %q", b.namespace, "default")
	}
	if b.socketPath != "/run/containerd/containerd.sock" {
		t.Errorf("socketPath = %q, want %q", b.socketPath, "/run/containerd/containerd.sock")
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
