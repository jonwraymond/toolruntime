package remote

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
)

func TestBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*Backend)(nil)
}

func TestBackendKind(t *testing.T) {
	b := New(Config{Endpoint: "http://localhost:8080"})
	if b.Kind() != toolruntime.BackendRemote {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendRemote)
	}
}

func TestBackendDefaults(t *testing.T) {
	b := New(Config{Endpoint: "http://localhost:8080"})
	if b.timeoutOverhead != 5*time.Second {
		t.Errorf("timeoutOverhead = %v, want %v", b.timeoutOverhead, 5*time.Second)
	}
	if b.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want %d", b.maxRetries, 3)
	}
}

func TestBackendRequiresGateway(t *testing.T) {
	b := New(Config{Endpoint: "http://localhost:8080"})
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

func TestBackendRequiresEndpoint(t *testing.T) {
	b := New(Config{Endpoint: ""}) // No endpoint
	ctx := context.Background()

	// Create a mock gateway
	gw := &mockGateway{}
	req := toolruntime.ExecuteRequest{
		Code:    "test",
		Gateway: gw,
	}
	_, err := b.Execute(ctx, req)
	if !errors.Is(err, ErrRemoteNotAvailable) {
		t.Errorf("Execute() without endpoint error = %v, want %v", err, ErrRemoteNotAvailable)
	}
}

// mockGateway implements toolruntime.ToolGateway for testing
type mockGateway struct{}

func (m *mockGateway) SearchTools(_ context.Context, _ string, _ int) ([]toolindex.Summary, error) {
	return nil, nil
}
func (m *mockGateway) ListNamespaces(_ context.Context) ([]string, error) {
	return nil, nil
}
func (m *mockGateway) DescribeTool(_ context.Context, _ string, _ tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, nil
}
func (m *mockGateway) ListToolExamples(_ context.Context, _ string, _ int) ([]tooldocs.ToolExample, error) {
	return nil, nil
}
func (m *mockGateway) RunTool(_ context.Context, _ string, _ map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, nil
}
func (m *mockGateway) RunChain(_ context.Context, _ []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, nil
}
