package toolcodeengine

import (
	"context"

	"github.com/jonwraymond/toolcode"
	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
)

// toolsGateway wraps toolcode.Tools to implement toolruntime.ToolGateway.
type toolsGateway struct {
	tools toolcode.Tools
}

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// WrapTools wraps a toolcode.Tools implementation to satisfy toolruntime.ToolGateway.
// This allows the toolcode.Tools interface to be used as a gateway in toolruntime.
func WrapTools(tools toolcode.Tools) toolruntime.ToolGateway {
	return &toolsGateway{tools: tools}
}

// SearchTools implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return g.tools.SearchTools(ctx, query, limit)
}

// ListNamespaces implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return g.tools.ListNamespaces(ctx)
}

// DescribeTool implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return tooldocs.ToolDoc{}, err
	}
	return g.tools.DescribeTool(ctx, id, level)
}

// ListToolExamples implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) ListToolExamples(ctx context.Context, id string, maxExamples int) ([]tooldocs.ToolExample, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return g.tools.ListToolExamples(ctx, id, maxExamples)
}

// RunTool implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return toolrun.RunResult{}, err
	}
	return g.tools.RunTool(ctx, id, args)
}

// RunChain implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	ctx = ensureContext(ctx)
	if err := ctx.Err(); err != nil {
		return toolrun.RunResult{}, nil, err
	}
	return g.tools.RunChain(ctx, steps)
}
