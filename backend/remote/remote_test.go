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

func TestRemoteBackendImplementsInterface(t *testing.T) {
	var _ toolruntime.Backend = (*RemoteBackend)(nil)
}

func TestRemoteBackendKind(t *testing.T) {
	b := New(Config{Endpoint: "http://localhost:8080"})
	if b.Kind() != toolruntime.BackendRemote {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendRemote)
	}
}

func TestRemoteBackendDefaults(t *testing.T) {
	b := New(Config{Endpoint: "http://localhost:8080"})
	if b.timeoutOverhead != 5*time.Second {
		t.Errorf("timeoutOverhead = %v, want %v", b.timeoutOverhead, 5*time.Second)
	}
	if b.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want %d", b.maxRetries, 3)
	}
}

func TestRemoteBackendRequiresGateway(t *testing.T) {
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

func TestRemoteBackendRequiresEndpoint(t *testing.T) {
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

func (m *mockGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	return nil, nil
}
func (m *mockGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (m *mockGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, nil
}
func (m *mockGateway) ListToolExamples(ctx context.Context, id string, max int) ([]tooldocs.ToolExample, error) {
	return nil, nil
}
func (m *mockGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, nil
}
func (m *mockGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, nil
}
