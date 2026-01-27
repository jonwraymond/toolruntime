package toolruntime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
)

// Test that mockToolGateway satisfies the interface
func TestMockToolGatewayImplementsInterface(t *testing.T) {
	var _ ToolGateway = (*mockToolGateway)(nil)
}

// Test mock gateway contract compliance
func TestMockGatewayContract(t *testing.T) {
	RunGatewayContractTests(t, GatewayContract{
		NewGateway: func() ToolGateway {
			return &mockToolGateway{}
		},
	})
}

func TestRecordingGatewayCounts(t *testing.T) {
	base := &mockToolGateway{}
	rec := &recordingGateway{wrapped: base}

	ctx := context.Background()
	if _, err := rec.SearchTools(ctx, "q", 1); err != nil {
		t.Fatalf("SearchTools failed: %v", err)
	}
	if _, err := rec.RunTool(ctx, "tool", map[string]any{"k": "v"}); err != nil {
		t.Fatalf("RunTool failed: %v", err)
	}
	if _, _, err := rec.RunChain(ctx, []toolrun.ChainStep{{ToolID: "tool"}}); err != nil {
		t.Fatalf("RunChain failed: %v", err)
	}

	if rec.searchCalls != 1 {
		t.Fatalf("searchCalls = %d, want 1", rec.searchCalls)
	}
	if rec.runToolCalls != 1 {
		t.Fatalf("runToolCalls = %d, want 1", rec.runToolCalls)
	}
	if rec.runChainCalls != 1 {
		t.Fatalf("runChainCalls = %d, want 1", rec.runChainCalls)
	}
}

func TestErrGatewayPropagates(t *testing.T) {
	expected := errors.New("boom")
	g := &errGateway{err: expected}
	ctx := context.Background()

	if _, err := g.SearchTools(ctx, "q", 1); !errors.Is(err, expected) {
		t.Fatalf("SearchTools error = %v, want %v", err, expected)
	}
	if _, err := g.ListNamespaces(ctx); !errors.Is(err, expected) {
		t.Fatalf("ListNamespaces error = %v, want %v", err, expected)
	}
	if _, err := g.DescribeTool(ctx, "id", tooldocs.DetailSummary); !errors.Is(err, expected) {
		t.Fatalf("DescribeTool error = %v, want %v", err, expected)
	}
	if _, err := g.ListToolExamples(ctx, "id", 1); !errors.Is(err, expected) {
		t.Fatalf("ListToolExamples error = %v, want %v", err, expected)
	}
	if _, err := g.RunTool(ctx, "id", nil); !errors.Is(err, expected) {
		t.Fatalf("RunTool error = %v, want %v", err, expected)
	}
	if _, _, err := g.RunChain(ctx, []toolrun.ChainStep{{ToolID: "id"}}); !errors.Is(err, expected) {
		t.Fatalf("RunChain error = %v, want %v", err, expected)
	}
}

// recordingGateway wraps a gateway and records calls for testing
type recordingGateway struct {
	wrapped       ToolGateway
	searchCalls   int
	runToolCalls  int
	runChainCalls int
}

func (r *recordingGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	r.searchCalls++
	return r.wrapped.SearchTools(ctx, query, limit)
}

func (r *recordingGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return r.wrapped.ListNamespaces(ctx)
}

func (r *recordingGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return r.wrapped.DescribeTool(ctx, id, level)
}

func (r *recordingGateway) ListToolExamples(ctx context.Context, id string, max int) ([]tooldocs.ToolExample, error) {
	return r.wrapped.ListToolExamples(ctx, id, max)
}

func (r *recordingGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	r.runToolCalls++
	return r.wrapped.RunTool(ctx, id, args)
}

func (r *recordingGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	r.runChainCalls++
	return r.wrapped.RunChain(ctx, steps)
}

// errGateway returns errors for testing error handling
type errGateway struct {
	err error
}

func (e *errGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	return nil, e.err
}

func (e *errGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return nil, e.err
}

func (e *errGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, e.err
}

func (e *errGateway) ListToolExamples(ctx context.Context, id string, max int) ([]tooldocs.ToolExample, error) {
	return nil, e.err
}

func (e *errGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, e.err
}

func (e *errGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, e.err
}

// delayGateway adds delays for testing timeouts
type delayGateway struct {
	wrapped ToolGateway
	delay   time.Duration
}

func (d *delayGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.SearchTools(ctx, query, limit)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *delayGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.ListNamespaces(ctx)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *delayGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.DescribeTool(ctx, id, level)
	case <-ctx.Done():
		return tooldocs.ToolDoc{}, ctx.Err()
	}
}

func (d *delayGateway) ListToolExamples(ctx context.Context, id string, max int) ([]tooldocs.ToolExample, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.ListToolExamples(ctx, id, max)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *delayGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.RunTool(ctx, id, args)
	case <-ctx.Done():
		return toolrun.RunResult{}, ctx.Err()
	}
}

func (d *delayGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	select {
	case <-time.After(d.delay):
		return d.wrapped.RunChain(ctx, steps)
	case <-ctx.Done():
		return toolrun.RunResult{}, nil, ctx.Err()
	}
}

func TestDelayGatewayTimeout(t *testing.T) {
	base := &mockToolGateway{}
	delayed := &delayGateway{wrapped: base, delay: 100 * time.Millisecond}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := delayed.SearchTools(ctx, "test", 10)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
