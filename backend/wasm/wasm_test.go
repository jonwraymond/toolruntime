package wasm

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/toolruntime"
)

func TestWASMBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*WASMBackend)(nil)
}

func TestWASMBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendWASM {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendWASM)
	}
}

func TestWASMBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.runtime != "wazero" {
		t.Errorf("runtime = %q, want %q", b.runtime, "wazero")
	}
	if b.maxMemoryPages != 256 {
		t.Errorf("maxMemoryPages = %d, want %d", b.maxMemoryPages, 256)
	}
}

func TestWASMBackendRequiresGateway(t *testing.T) {
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
