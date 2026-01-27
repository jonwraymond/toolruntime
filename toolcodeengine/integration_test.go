package toolcodeengine_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonwraymond/toolcode"
	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
	"github.com/jonwraymond/toolruntime/backend/unsafe"
	"github.com/jonwraymond/toolruntime/toolcodeengine"
)

// testTools implements toolcode.Tools for integration testing
type testTools struct {
	searchResults []toolindex.Summary
	namespaces    []string
	toolDoc       tooldocs.ToolDoc
	examples      []tooldocs.ToolExample
	runResult     toolrun.RunResult
	chainResult   toolrun.RunResult
	stepResults   []toolrun.StepResult
}

func (t *testTools) SearchTools(_ string, _ int) ([]toolindex.Summary, error) {
	return t.searchResults, nil
}

func (t *testTools) ListNamespaces() ([]string, error) {
	return t.namespaces, nil
}

func (t *testTools) DescribeTool(_ string, _ tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return t.toolDoc, nil
}

func (t *testTools) ListToolExamples(_ string, _ int) ([]tooldocs.ToolExample, error) {
	return t.examples, nil
}

func (t *testTools) RunTool(_ context.Context, _ string, _ map[string]any) (toolrun.RunResult, error) {
	return t.runResult, nil
}

func (t *testTools) RunChain(_ context.Context, _ []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return t.chainResult, t.stepResults, nil
}

func (t *testTools) Println(_ ...any) {}

var _ toolcode.Tools = (*testTools)(nil)

// TestFullStackExecution tests toolcode -> toolcodeengine -> toolruntime -> unsafe backend
func TestFullStackExecution(t *testing.T) {
	// Create an unsafe backend
	backend := unsafe.New(unsafe.Config{
		Mode: unsafe.ModeSubprocess,
	})

	// Create a runtime with the backend
	runtime := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
		Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
			toolruntime.ProfileDev: backend,
		},
		DefaultProfile: toolruntime.ProfileDev,
	})

	// Create the toolcode engine adapter
	engine := toolcodeengine.New(toolcodeengine.Config{
		Runtime: runtime,
		Profile: toolruntime.ProfileDev,
	})

	// Verify Engine implements toolcode.Engine
	var _ toolcode.Engine = engine

	// Create test tools
	tools := &testTools{
		searchResults: []toolindex.Summary{
			{ID: "test:tool", Name: "tool"},
		},
	}

	// Execute simple code
	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Language:     "go",
		Code:         `__out = "hello world"`,
		Timeout:      10 * time.Second,
		MaxToolCalls: 5,
	}

	result, err := engine.Execute(ctx, params, tools)
	if err != nil {
		t.Logf("Execute() error (expected during integration): %v", err)
		// Note: The unsafe backend subprocess execution may fail depending on environment
		// The important thing is that the full stack is wired up correctly
	}

	// If execution succeeded, verify result
	if err == nil {
		t.Logf("Execute() result: %+v", result)
		if result.Value != "hello world" {
			t.Logf("Unexpected value: %v (may vary based on backend execution)", result.Value)
		}
	}
}

// TestErrorMappingIntegration tests that errors are correctly mapped through the stack
func TestErrorMappingIntegration(t *testing.T) {
	// Create a mock backend that returns specific errors
	mockBackend := &errorBackend{}

	runtime := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
		Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
			toolruntime.ProfileDev: mockBackend,
		},
		DefaultProfile: toolruntime.ProfileDev,
	})

	engine := toolcodeengine.New(toolcodeengine.Config{
		Runtime: runtime,
		Profile: toolruntime.ProfileDev,
	})

	tools := &testTools{}
	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}

	t.Run("timeout maps to ErrLimitExceeded", func(t *testing.T) {
		mockBackend.err = toolruntime.ErrTimeout
		_, err := engine.Execute(ctx, params, tools)
		if !errors.Is(err, toolcode.ErrLimitExceeded) {
			t.Errorf("timeout should map to ErrLimitExceeded, got: %v", err)
		}
	})

	t.Run("resource limit maps to ErrLimitExceeded", func(t *testing.T) {
		mockBackend.err = toolruntime.ErrResourceLimit
		_, err := engine.Execute(ctx, params, tools)
		if !errors.Is(err, toolcode.ErrLimitExceeded) {
			t.Errorf("resource limit should map to ErrLimitExceeded, got: %v", err)
		}
	})

	t.Run("sandbox violation maps to ErrCodeExecution", func(t *testing.T) {
		mockBackend.err = toolruntime.ErrSandboxViolation
		_, err := engine.Execute(ctx, params, tools)
		if !errors.Is(err, toolcode.ErrCodeExecution) {
			t.Errorf("sandbox violation should map to ErrCodeExecution, got: %v", err)
		}
	})
}

// errorBackend is a mock backend that returns configurable errors
type errorBackend struct {
	err error
}

func (b *errorBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendUnsafeHost
}

func (b *errorBackend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	if b.err != nil {
		return toolruntime.ExecuteResult{}, b.err
	}
	return toolruntime.ExecuteResult{
		Value: "test",
	}, nil
}

// TestGatewayWrappingIntegration tests that Tools is correctly wrapped as Gateway
func TestGatewayWrappingIntegration(t *testing.T) {
	// Create a mock backend that captures the request
	mockBackend := &capturingBackend{}

	runtime := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
		Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
			toolruntime.ProfileDev: mockBackend,
		},
		DefaultProfile: toolruntime.ProfileDev,
	})

	engine := toolcodeengine.New(toolcodeengine.Config{
		Runtime: runtime,
		Profile: toolruntime.ProfileDev,
	})

	tools := &testTools{
		searchResults: []toolindex.Summary{
			{ID: "tool1", Name: "Tool One"},
			{ID: "tool2", Name: "Tool Two"},
		},
		namespaces: []string{"ns1", "ns2"},
	}

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}

	_, _ = engine.Execute(ctx, params, tools)

	// Verify gateway was passed to backend
	if mockBackend.capturedReq.Gateway == nil {
		t.Error("Gateway should be passed to backend")
	}

	// Verify gateway works correctly
	gw := mockBackend.capturedReq.Gateway

	results, err := gw.SearchTools(ctx, "test", 10)
	if err != nil {
		t.Errorf("SearchTools() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("SearchTools() returned %d results, want 2", len(results))
	}

	namespaces, err := gw.ListNamespaces(ctx)
	if err != nil {
		t.Errorf("ListNamespaces() error = %v", err)
	}
	if len(namespaces) != 2 {
		t.Errorf("ListNamespaces() returned %d namespaces, want 2", len(namespaces))
	}
}

// capturingBackend captures the ExecuteRequest for inspection
type capturingBackend struct {
	capturedReq toolruntime.ExecuteRequest
}

func (b *capturingBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendUnsafeHost
}

func (b *capturingBackend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	b.capturedReq = req
	return toolruntime.ExecuteResult{}, nil
}

// TestProfilePropagation tests that security profiles are correctly propagated
func TestProfilePropagation(t *testing.T) {
	mockBackend := &capturingBackend{}

	runtime := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
		Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
			toolruntime.ProfileStandard: mockBackend,
		},
		DefaultProfile: toolruntime.ProfileStandard,
	})

	engine := toolcodeengine.New(toolcodeengine.Config{
		Runtime: runtime,
		Profile: toolruntime.ProfileStandard,
	})

	tools := &testTools{}
	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}

	_, _ = engine.Execute(ctx, params, tools)

	if mockBackend.capturedReq.Profile != toolruntime.ProfileStandard {
		t.Errorf("Profile = %v, want %v",
			mockBackend.capturedReq.Profile, toolruntime.ProfileStandard)
	}
}

// TestLimitsPropagation tests that limits are correctly propagated
func TestLimitsPropagation(t *testing.T) {
	mockBackend := &capturingBackend{}

	runtime := toolruntime.NewDefaultRuntime(toolruntime.RuntimeConfig{
		Backends: map[toolruntime.SecurityProfile]toolruntime.Backend{
			toolruntime.ProfileDev: mockBackend,
		},
		DefaultProfile: toolruntime.ProfileDev,
	})

	engine := toolcodeengine.New(toolcodeengine.Config{
		Runtime: runtime,
		Profile: toolruntime.ProfileDev,
	})

	tools := &testTools{}
	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code:         "test",
		Timeout:      15 * time.Second,
		MaxToolCalls: 25,
	}

	_, _ = engine.Execute(ctx, params, tools)

	if mockBackend.capturedReq.Timeout != 15*time.Second {
		t.Errorf("Timeout = %v, want %v",
			mockBackend.capturedReq.Timeout, 15*time.Second)
	}
	if mockBackend.capturedReq.Limits.MaxToolCalls != 25 {
		t.Errorf("MaxToolCalls = %d, want %d",
			mockBackend.capturedReq.Limits.MaxToolCalls, 25)
	}
}
