package temporal

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
)

// mockBackend implements toolruntime.Backend for testing
type mockBackend struct{}

func (m *mockBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendDocker
}

func (m *mockBackend) Execute(_ context.Context, _ toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	return toolruntime.ExecuteResult{}, nil
}

func TestBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*Backend)(nil)
}

func TestBackendKind(t *testing.T) {
	b := New(Config{})
	if b.Kind() != toolruntime.BackendTemporal {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendTemporal)
	}
}

func TestBackendDefaults(t *testing.T) {
	b := New(Config{})
	if b.hostPort != "localhost:7233" {
		t.Errorf("hostPort = %q, want %q", b.hostPort, "localhost:7233")
	}
	if b.namespace != "default" {
		t.Errorf("namespace = %q, want %q", b.namespace, "default")
	}
	if b.taskQueue != "toolruntime-execution" {
		t.Errorf("taskQueue = %q, want %q", b.taskQueue, "toolruntime-execution")
	}
}

func TestBackendRequiresGateway(t *testing.T) {
	b := New(Config{
		SandboxBackend: &mockBackend{},
	})
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

func TestBackendRequiresSandboxBackend(t *testing.T) {
	b := New(Config{
		SandboxBackend: nil, // No sandbox backend
	})
	ctx := context.Background()

	// Create a mock gateway
	gw := &mockGateway{}
	req := toolruntime.ExecuteRequest{
		Code:    "test",
		Gateway: gw,
	}
	_, err := b.Execute(ctx, req)
	if !errors.Is(err, ErrMissingSandboxBackend) {
		t.Errorf("Execute() without sandbox backend error = %v, want %v", err, ErrMissingSandboxBackend)
	}
}

// mockGateway implements toolruntime.ToolGateway for testing
type mockGateway struct{}

func (m *mockGateway) SearchTools(ctx context.Context, _ string, _ int) ([]toolindex.Summary, error) {
	return nil, nil
}
func (m *mockGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (m *mockGateway) DescribeTool(ctx context.Context, _ string, _ tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, nil
}
func (m *mockGateway) ListToolExamples(ctx context.Context, _ string, _ int) ([]tooldocs.ToolExample, error) {
	return nil, nil
}
func (m *mockGateway) RunTool(ctx context.Context, _ string, _ map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, nil
}
func (m *mockGateway) RunChain(ctx context.Context, _ []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, nil
}
