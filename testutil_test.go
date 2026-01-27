package toolruntime

import (
	"context"
	"errors"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
)

// errToolNotFound is used by mock gateway
var errToolNotFound = errors.New("tool not found")

// mockToolGateway is a minimal mock for testing
type mockToolGateway struct{}

func (m *mockToolGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

func (m *mockToolGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

func (m *mockToolGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	if ctx.Err() != nil {
		return tooldocs.ToolDoc{}, ctx.Err()
	}
	// Return error for non-existent tools (anything we don't know about)
	return tooldocs.ToolDoc{}, errToolNotFound
}

func (m *mockToolGateway) ListToolExamples(ctx context.Context, id string, max int) ([]tooldocs.ToolExample, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

func (m *mockToolGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	if ctx.Err() != nil {
		return toolrun.RunResult{}, ctx.Err()
	}
	return toolrun.RunResult{}, nil
}

func (m *mockToolGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	if ctx.Err() != nil {
		return toolrun.RunResult{}, nil, ctx.Err()
	}
	return toolrun.RunResult{}, nil, nil
}
