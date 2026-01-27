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

// WrapTools wraps a toolcode.Tools implementation to satisfy toolruntime.ToolGateway.
// This allows the toolcode.Tools interface to be used as a gateway in toolruntime.
func WrapTools(tools toolcode.Tools) toolruntime.ToolGateway {
	return &toolsGateway{tools: tools}
}

// SearchTools implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) SearchTools(ctx context.Context, query string, limit int) ([]toolindex.Summary, error) {
	return g.tools.SearchTools(query, limit)
}

// ListNamespaces implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return g.tools.ListNamespaces()
}

// DescribeTool implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) DescribeTool(ctx context.Context, id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return g.tools.DescribeTool(id, level)
}

// ListToolExamples implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) ListToolExamples(ctx context.Context, id string, maxExamples int) ([]tooldocs.ToolExample, error) {
	return g.tools.ListToolExamples(id, maxExamples)
}

// RunTool implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	return g.tools.RunTool(ctx, id, args)
}

// RunChain implements toolruntime.ToolGateway by delegating to the wrapped Tools.
func (g *toolsGateway) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return g.tools.RunChain(ctx, steps)
}
