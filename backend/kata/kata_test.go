package kata

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/toolruntime"
)

func TestKataBackendImplementsInterface(t *testing.T) {
	var _ toolruntime.Backend = (*KataBackend)(nil)
}

func TestKataBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendKata {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendKata)
	}
}

func TestKataBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.runtimePath != "kata-runtime" {
		t.Errorf("runtimePath = %q, want %q", b.runtimePath, "kata-runtime")
	}
	if b.hypervisor != "qemu" {
		t.Errorf("hypervisor = %q, want %q", b.hypervisor, "qemu")
	}
}

func TestKataBackendRequiresGateway(t *testing.T) {
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
