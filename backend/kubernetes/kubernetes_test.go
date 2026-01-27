package kubernetes

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/toolruntime"
)

func TestKubernetesBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*KubernetesBackend)(nil)
}

func TestKubernetesBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendKubernetes {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendKubernetes)
	}
}

func TestKubernetesBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.namespace != "default" {
		t.Errorf("namespace = %q, want %q", b.namespace, "default")
	}
	if b.image != "toolruntime-sandbox:latest" {
		t.Errorf("image = %q, want %q", b.image, "toolruntime-sandbox:latest")
	}
}

func TestKubernetesBackendRequiresGateway(t *testing.T) {
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
