package toolcodeengine

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
)

// mockRuntime implements toolruntime.Runtime for testing
type mockRuntime struct {
	result      toolruntime.ExecuteResult
	err         error
	capturedReq toolruntime.ExecuteRequest
}

func (m *mockRuntime) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	m.capturedReq = req
	if m.err != nil {
		return toolruntime.ExecuteResult{}, m.err
	}
	return m.result, nil
}

// mockTools implements toolcode.Tools for testing
type mockTools struct {
	searchResults []toolindex.Summary
	namespaces    []string
	toolDoc       tooldocs.ToolDoc
	examples      []tooldocs.ToolExample
	runResult     toolrun.RunResult
	chainResult   toolrun.RunResult
	stepResults   []toolrun.StepResult
}

func (m *mockTools) SearchTools(query string, limit int) ([]toolindex.Summary, error) {
	return m.searchResults, nil
}

func (m *mockTools) ListNamespaces() ([]string, error) {
	return m.namespaces, nil
}

func (m *mockTools) DescribeTool(id string, level tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return m.toolDoc, nil
}

func (m *mockTools) ListToolExamples(id string, max int) ([]tooldocs.ToolExample, error) {
	return m.examples, nil
}

func (m *mockTools) RunTool(ctx context.Context, id string, args map[string]any) (toolrun.RunResult, error) {
	return m.runResult, nil
}

func (m *mockTools) RunChain(ctx context.Context, steps []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return m.chainResult, m.stepResults, nil
}

func (m *mockTools) Println(args ...any) {
	// Mock implementation
}

// TestEngineImplementsInterface verifies Engine satisfies toolcode.Engine
func TestEngineImplementsInterface(t *testing.T) {
	var _ toolcode.Engine = (*Engine)(nil)
}

func TestEngineExecute(t *testing.T) {
	rt := &mockRuntime{
		result: toolruntime.ExecuteResult{
			Value:  "result",
			Stdout: "output",
			Stderr: "",
			Duration: 100 * time.Millisecond,
		},
	}

	engine := New(Config{
		Runtime: rt,
		Profile: toolruntime.ProfileDev,
	})

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Language:     "go",
		Code:         `__out = "hello"`,
		Timeout:      30 * time.Second,
		MaxToolCalls: 10,
	}
	tools := &mockTools{}

	result, err := engine.Execute(ctx, params, tools)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Value != "result" {
		t.Errorf("Execute().Value = %v, want %v", result.Value, "result")
	}
	if result.Stdout != "output" {
		t.Errorf("Execute().Stdout = %q, want %q", result.Stdout, "output")
	}
}

func TestEngineExecuteMapsParams(t *testing.T) {
	rt := &mockRuntime{
		result: toolruntime.ExecuteResult{},
	}

	engine := New(Config{
		Runtime: rt,
		Profile: toolruntime.ProfileStandard,
	})

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Language:     "go",
		Code:         "test code",
		Timeout:      5 * time.Second,
		MaxToolCalls: 20,
	}
	tools := &mockTools{}

	_, _ = engine.Execute(ctx, params, tools)

	if rt.capturedReq.Language != "go" {
		t.Errorf("ExecuteRequest.Language = %q, want %q", rt.capturedReq.Language, "go")
	}
	if rt.capturedReq.Code != "test code" {
		t.Errorf("ExecuteRequest.Code = %q, want %q", rt.capturedReq.Code, "test code")
	}
	if rt.capturedReq.Timeout != 5*time.Second {
		t.Errorf("ExecuteRequest.Timeout = %v, want %v", rt.capturedReq.Timeout, 5*time.Second)
	}
	if rt.capturedReq.Limits.MaxToolCalls != 20 {
		t.Errorf("ExecuteRequest.Limits.MaxToolCalls = %d, want %d", rt.capturedReq.Limits.MaxToolCalls, 20)
	}
	if rt.capturedReq.Profile != toolruntime.ProfileStandard {
		t.Errorf("ExecuteRequest.Profile = %v, want %v", rt.capturedReq.Profile, toolruntime.ProfileStandard)
	}
}

func TestEngineExecuteTimeoutError(t *testing.T) {
	rt := &mockRuntime{
		err: toolruntime.ErrTimeout,
	}

	engine := New(Config{
		Runtime: rt,
		Profile: toolruntime.ProfileDev,
	})

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}
	tools := &mockTools{}

	_, err := engine.Execute(ctx, params, tools)
	if err == nil {
		t.Error("Execute() should return error for timeout")
	}
	// Should map to toolcode.ErrLimitExceeded
	if !errors.Is(err, toolcode.ErrLimitExceeded) {
		t.Errorf("Execute() error = %v, want wrapped %v", err, toolcode.ErrLimitExceeded)
	}
}

func TestEngineExecuteResourceLimitError(t *testing.T) {
	rt := &mockRuntime{
		err: toolruntime.ErrResourceLimit,
	}

	engine := New(Config{
		Runtime: rt,
		Profile: toolruntime.ProfileDev,
	})

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}
	tools := &mockTools{}

	_, err := engine.Execute(ctx, params, tools)
	if err == nil {
		t.Error("Execute() should return error for resource limit")
	}
	if !errors.Is(err, toolcode.ErrLimitExceeded) {
		t.Errorf("Execute() error = %v, want wrapped %v", err, toolcode.ErrLimitExceeded)
	}
}

func TestEngineExecuteSandboxViolationError(t *testing.T) {
	rt := &mockRuntime{
		err: toolruntime.ErrSandboxViolation,
	}

	engine := New(Config{
		Runtime: rt,
		Profile: toolruntime.ProfileDev,
	})

	ctx := context.Background()
	params := toolcode.ExecuteParams{
		Code: "test",
	}
	tools := &mockTools{}

	_, err := engine.Execute(ctx, params, tools)
	if err == nil {
		t.Error("Execute() should return error for sandbox violation")
	}
	if !errors.Is(err, toolcode.ErrCodeExecution) {
		t.Errorf("Execute() error = %v, want wrapped %v", err, toolcode.ErrCodeExecution)
	}
}
