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
