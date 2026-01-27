package direct

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolmodel"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
)

// mockIndex implements toolindex.Index for testing
type mockIndex struct {
	summaries  []toolindex.Summary
	namespaces []string
	searchErr  error
}

func (m *mockIndex) RegisterTool(tool toolmodel.Tool, backend toolmodel.ToolBackend) error {
	return nil
}

func (m *mockIndex) RegisterTools(regs []toolindex.ToolRegistration) error {
	return nil
}

func (m *mockIndex) RegisterToolsFromMCP(serverName string, tools []toolmodel.Tool) error {
	return nil
}

func (m *mockIndex) UnregisterBackend(toolID string, kind toolmodel.BackendKind, backendID string) error {
	return nil
}

func (m *mockIndex) GetTool(id string) (toolmodel.Tool, toolmodel.ToolBackend, error) {
	return toolmodel.Tool{}, toolmodel.ToolBackend{}, nil
}

func (m *mockIndex) GetAllBackends(id string) ([]toolmodel.ToolBackend, error) {
	return nil, nil
}

func (m *mockIndex) Search(query string, limit int) ([]toolindex.Summary, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	if limit > len(m.summaries) {
		return m.summaries, nil
	}
	return m.summaries[:limit], nil
}

func (m *mockIndex) ListNamespaces() ([]string, error) {
	return m.namespaces, nil
}

// mockDocs implements tooldocs.Store for testing
type mockDocs struct {
	docs       map[string]tooldocs.ToolDoc
	examples   map[string][]tooldocs.ToolExample
	descErr    error
	examplesErr error
}

func (m *mockDocs) DescribeTool(id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	if m.descErr != nil {
		return tooldocs.ToolDoc{}, m.descErr
	}
	doc, ok := m.docs[id]
	if !ok {
		return tooldocs.ToolDoc{}, errors.New("tool not found")
	}
	return doc, nil
}

func (m *mockDocs) ListExamples(id string, max int) ([]tooldocs.ToolExample, error) {
	if m.examplesErr != nil {
		return nil, m.examplesErr
	}
	examples := m.examples[id]
	if max > 0 && max < len(examples) {
		return examples[:max], nil
	}
	return examples, nil
}

// mockRunner implements toolrun.Runner for testing
type mockRunner struct {
	runResult    toolrun.RunResult
	runErr       error
	chainResult  toolrun.RunResult
	stepResults  []toolrun.StepResult
	chainErr     error
	callCount    int
	mu           sync.Mutex
}

func (m *mockRunner) Run(ctx context.Context, toolID string, args map[string]any) (toolrun.RunResult, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	if ctx.Err() != nil {
		return toolrun.RunResult{}, ctx.Err()
	}
	if m.runErr != nil {
		return toolrun.RunResult{}, m.runErr
	}
	return m.runResult, nil
}

func (m *mockRunner) RunStream(ctx context.Context, toolID string, args map[string]any) (<-chan toolrun.StreamEvent, error) {
	return nil, errors.New("streaming not supported")
}

func (m *mockRunner) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	m.mu.Lock()
	m.callCount += len(steps)
	m.mu.Unlock()

	if ctx.Err() != nil {
		return toolrun.RunResult{}, nil, ctx.Err()
	}
	if m.chainErr != nil {
		return toolrun.RunResult{}, nil, m.chainErr
	}
	return m.chainResult, m.stepResults, nil
}

// TestDirectGatewayImplementsInterface verifies DirectGateway satisfies ToolGateway
func TestDirectGatewayImplementsInterface(t *testing.T) {
	var _ toolruntime.ToolGateway = (*DirectGateway)(nil)
}

func TestDirectGatewaySearchTools(t *testing.T) {
	summaries := []toolindex.Summary{
		{ID: "test:tool1", Name: "tool1"},
		{ID: "test:tool2", Name: "tool2"},
		{ID: "test:tool3", Name: "tool3"},
	}

	index := &mockIndex{summaries: summaries}
	gw := New(Config{
		Index: index,
		Docs:  &mockDocs{},
		Runner: &mockRunner{},
	})

	ctx := context.Background()

	t.Run("delegates to index", func(t *testing.T) {
		results, err := gw.SearchTools(ctx, "test", 10)
		if err != nil {
			t.Fatalf("SearchTools() error = %v", err)
		}
		if len(results) != 3 {
			t.Errorf("SearchTools() returned %d results, want 3", len(results))
		}
	})

	t.Run("respects limit", func(t *testing.T) {
		results, err := gw.SearchTools(ctx, "test", 2)
		if err != nil {
			t.Fatalf("SearchTools() error = %v", err)
		}
		if len(results) != 2 {
			t.Errorf("SearchTools() returned %d results, want 2", len(results))
		}
	})

	t.Run("propagates errors", func(t *testing.T) {
		index := &mockIndex{searchErr: errors.New("search failed")}
		gw := New(Config{
			Index: index,
			Docs:  &mockDocs{},
			Runner: &mockRunner{},
		})

		_, err := gw.SearchTools(ctx, "test", 10)
		if err == nil {
			t.Error("SearchTools() should propagate error")
		}
	})
}

func TestDirectGatewayListNamespaces(t *testing.T) {
	namespaces := []string{"ns1", "ns2"}
	index := &mockIndex{namespaces: namespaces}
	gw := New(Config{
		Index: index,
		Docs:  &mockDocs{},
		Runner: &mockRunner{},
	})

	ctx := context.Background()

	results, err := gw.ListNamespaces(ctx)
	if err != nil {
		t.Fatalf("ListNamespaces() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("ListNamespaces() returned %d results, want 2", len(results))
	}
}

func TestDirectGatewayDescribeTool(t *testing.T) {
	docs := &mockDocs{
		docs: map[string]tooldocs.ToolDoc{
			"test:tool": {Summary: "Test tool"},
		},
	}
	gw := New(Config{
		Index: &mockIndex{},
		Docs:  docs,
		Runner: &mockRunner{},
	})

	ctx := context.Background()

	t.Run("delegates to docs store", func(t *testing.T) {
		doc, err := gw.DescribeTool(ctx, "test:tool", tooldocs.DetailSummary)
		if err != nil {
			t.Fatalf("DescribeTool() error = %v", err)
		}
		if doc.Summary != "Test tool" {
			t.Errorf("DescribeTool().Summary = %q, want %q", doc.Summary, "Test tool")
		}
	})

	t.Run("returns error for non-existent tool", func(t *testing.T) {
		_, err := gw.DescribeTool(ctx, "nonexistent:tool", tooldocs.DetailSummary)
		if err == nil {
			t.Error("DescribeTool() should return error for non-existent tool")
		}
	})
}

func TestDirectGatewayListToolExamples(t *testing.T) {
	examples := []tooldocs.ToolExample{
		{Title: "Example 1"},
		{Title: "Example 2"},
	}
	docs := &mockDocs{
		examples: map[string][]tooldocs.ToolExample{
			"test:tool": examples,
		},
	}
	gw := New(Config{
		Index: &mockIndex{},
		Docs:  docs,
		Runner: &mockRunner{},
	})

	ctx := context.Background()

	results, err := gw.ListToolExamples(ctx, "test:tool", 10)
	if err != nil {
		t.Fatalf("ListToolExamples() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("ListToolExamples() returned %d results, want 2", len(results))
	}
}

func TestDirectGatewayRunTool(t *testing.T) {
	runner := &mockRunner{
		runResult: toolrun.RunResult{
			Structured: "result",
		},
	}
	gw := New(Config{
		Index: &mockIndex{},
		Docs:  &mockDocs{},
		Runner: runner,
	})

	ctx := context.Background()

	t.Run("delegates to runner", func(t *testing.T) {
		result, err := gw.RunTool(ctx, "test:tool", map[string]any{"key": "value"})
		if err != nil {
			t.Fatalf("RunTool() error = %v", err)
		}
		if result.Structured != "result" {
			t.Errorf("RunTool().Structured = %v, want %v", result.Structured, "result")
		}
	})

	t.Run("propagates context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := gw.RunTool(ctx, "test:tool", nil)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("RunTool() should return context.Canceled, got %v", err)
		}
	})

	t.Run("records tool calls", func(t *testing.T) {
		runner := &mockRunner{}
		gw := New(Config{
			Index: &mockIndex{},
			Docs:  &mockDocs{},
			Runner: runner,
		})

		_, _ = gw.RunTool(context.Background(), "tool1", nil)
		_, _ = gw.RunTool(context.Background(), "tool2", nil)

		records := gw.GetToolCalls()
		if len(records) != 2 {
			t.Errorf("GetToolCalls() returned %d records, want 2", len(records))
		}
	})
}

func TestDirectGatewayRunChain(t *testing.T) {
	runner := &mockRunner{
		chainResult: toolrun.RunResult{
			Structured: "chain_result",
		},
		stepResults: []toolrun.StepResult{
			{ToolID: "step1"},
			{ToolID: "step2"},
		},
	}
	gw := New(Config{
		Index: &mockIndex{},
		Docs:  &mockDocs{},
		Runner: runner,
	})

	ctx := context.Background()

	t.Run("delegates to runner", func(t *testing.T) {
		steps := []toolrun.ChainStep{
			{ToolID: "step1"},
			{ToolID: "step2"},
		}
		result, stepResults, err := gw.RunChain(ctx, steps)
		if err != nil {
			t.Fatalf("RunChain() error = %v", err)
		}
		if result.Structured != "chain_result" {
			t.Errorf("RunChain().Structured = %v, want %v", result.Structured, "chain_result")
		}
		if len(stepResults) != 2 {
			t.Errorf("RunChain() returned %d step results, want 2", len(stepResults))
		}
	})

	t.Run("handles empty steps", func(t *testing.T) {
		_, stepResults, err := gw.RunChain(ctx, []toolrun.ChainStep{})
		if err != nil {
			t.Fatalf("RunChain() error = %v", err)
		}
		if len(stepResults) != 0 {
			t.Errorf("RunChain() with empty steps returned %d results", len(stepResults))
		}
	})
}

func TestDirectGatewayToolCallLimits(t *testing.T) {
	runner := &mockRunner{}
	gw := New(Config{
		Index:        &mockIndex{},
		Docs:         &mockDocs{},
		Runner:       runner,
		MaxToolCalls: 2,
	})

	ctx := context.Background()

	// First two calls should succeed
	_, err := gw.RunTool(ctx, "tool1", nil)
	if err != nil {
		t.Fatalf("RunTool() first call error = %v", err)
	}
	_, err = gw.RunTool(ctx, "tool2", nil)
	if err != nil {
		t.Fatalf("RunTool() second call error = %v", err)
	}

	// Third call should fail
	_, err = gw.RunTool(ctx, "tool3", nil)
	if err == nil {
		t.Error("RunTool() should fail after exceeding MaxToolCalls")
	}
	if !errors.Is(err, ErrToolCallLimitExceeded) {
		t.Errorf("RunTool() error = %v, want ErrToolCallLimitExceeded", err)
	}
}

func TestDirectGatewayChainStepLimits(t *testing.T) {
	runner := &mockRunner{
		stepResults: []toolrun.StepResult{{}, {}},
	}
	gw := New(Config{
		Index:         &mockIndex{},
		Docs:          &mockDocs{},
		Runner:        runner,
		MaxChainSteps: 2,
	})

	ctx := context.Background()

	// Chain with 3 steps should fail
	steps := []toolrun.ChainStep{
		{ToolID: "step1"},
		{ToolID: "step2"},
		{ToolID: "step3"},
	}
	_, _, err := gw.RunChain(ctx, steps)
	if err == nil {
		t.Error("RunChain() should fail when exceeding MaxChainSteps")
	}
	if !errors.Is(err, ErrChainStepLimitExceeded) {
		t.Errorf("RunChain() error = %v, want ErrChainStepLimitExceeded", err)
	}
}

func TestDirectGatewayThreadSafety(t *testing.T) {
	runner := &mockRunner{}
	gw := New(Config{
		Index:  &mockIndex{summaries: []toolindex.Summary{{ID: "test:tool"}}},
		Docs:   &mockDocs{},
		Runner: runner,
	})

	ctx := context.Background()
	var wg sync.WaitGroup

	// Run concurrent operations
	for i := 0; i < 10; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			_, _ = gw.SearchTools(ctx, "test", 10)
		}()
		go func() {
			defer wg.Done()
			_, _ = gw.RunTool(ctx, "test:tool", nil)
		}()
		go func() {
			defer wg.Done()
			_ = gw.GetToolCalls()
		}()
	}

	wg.Wait()
}

func TestDirectGatewayContractCompliance(t *testing.T) {
	toolruntime.RunGatewayContractTests(t, toolruntime.GatewayContract{
		NewGateway: func() toolruntime.ToolGateway {
			return New(Config{
				Index:  &mockIndex{summaries: []toolindex.Summary{{ID: "test:tool"}}},
				Docs:   &mockDocs{docs: map[string]tooldocs.ToolDoc{}},
				Runner: &mockRunner{},
			})
		},
	})
}
